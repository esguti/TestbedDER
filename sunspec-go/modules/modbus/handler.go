package modbus

import (
	"encoding/binary"

	"github.com/GoAethereal/cancel"
)

// Handler is firstly and foremost used by the modbus.Server.
// The Handle method describes how incoming messages are managed.
type Handler interface {
	Handle(ctx cancel.Context, code byte, req []byte) (res []byte, ex Exception)
}

var _ Handler = (*Mux)(nil)

// Mux implements the modbus.Handler interface and is intended to be used as a server side request
// multiplexer. When called by the server it will redirect the inbound message to the given function.
// If the callback is not set the Mux will return the modbus.ExIllegalFunction exception to the server.
// In case of an unknown function code the Fallback function, if set, will be executed.
// All given functions must be safe for use by multiple go routines.
type Mux struct {
	Fallback                   func(ctx cancel.Context, code byte, req []byte) (res []byte, ex Exception)
	ReadCoils                  func(ctx cancel.Context, address, quantity uint16) (res []bool, ex Exception)
	ReadDiscreteInputs         func(ctx cancel.Context, address, quantity uint16) (res []bool, ex Exception)
	ReadHoldingRegisters       func(ctx cancel.Context, address, quantity uint16) (res []byte, ex Exception)
	ReadInputRegisters         func(ctx cancel.Context, address, quantity uint16) (res []byte, ex Exception)
	WriteSingleCoil            func(ctx cancel.Context, address uint16, status bool) (ex Exception)
	WriteSingleRegister        func(ctx cancel.Context, address, value uint16) (ex Exception)
	WriteMultipleCoils         func(ctx cancel.Context, address uint16, status []bool) (ex Exception)
	WriteMultipleRegisters     func(ctx cancel.Context, address uint16, values []byte) (ex Exception)
	ReadWriteMultipleRegisters func(ctx cancel.Context, rAddress, rQuantity, wAddress uint16, values []byte) (res []byte, ex Exception)
}

// Handle dispatches incoming requests depending on their function code to the correlating callbacks
// as defined inside the Mux.
func (h *Mux) Handle(ctx cancel.Context, code byte, req []byte) (res []byte, ex Exception) {
	switch code {
	case 0x01:
		return h.readCoils(ctx, req)
	case 0x02:
		return h.readDiscreteInputs(ctx, req)
	case 0x03:
		return h.readHoldingRegisters(ctx, req)
	case 0x04:
		return h.readInputRegisters(ctx, req)
	case 0x05:
		return h.writeSingleCoil(ctx, req)
	case 0x06:
		return h.writeSingleRegister(ctx, req)
	case 0x0F:
		return h.writeMultipleCoils(ctx, req)
	case 0x10:
		return h.writeMultipleRegisters(ctx, req)
	case 0x17:
		return h.readWriteMultipleRegisters(ctx, req)
	}
	return h.fallback(ctx, code, req)
}

func (h *Mux) fallback(ctx cancel.Context, code byte, req []byte) (res []byte, ex Exception) {
	if h.Fallback == nil {
		return nil, IllegalFunction
	}
	return h.Fallback(ctx, code, req)
}

func (h *Mux) readCoils(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.ReadCoils == nil:
		return nil, IllegalFunction
	case len(req) != 4:
		return nil, IllegalDataAddress
	}
	address := binary.BigEndian.Uint16(req[0:])
	quantity := binary.BigEndian.Uint16(req[2:])
	if ex := boundCheck(address, quantity, 2000); ex != 0 {
		return nil, ex
	}
	status, ex := h.ReadCoils(ctx, address, quantity)
	switch {
	case ex != 0:
		return nil, ex
	case len(status) != int(quantity):
		return nil, SlaveDeviceFailure
	}
	return put(1+int(byteCount(quantity)), byte(byteCount(quantity)), status), 0
}

func (h *Mux) readDiscreteInputs(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.ReadDiscreteInputs == nil:
		return nil, IllegalFunction
	case len(req) != 4:
		return nil, IllegalDataAddress
	}
	address := binary.BigEndian.Uint16(req[0:])
	quantity := binary.BigEndian.Uint16(req[2:])
	if ex := boundCheck(address, quantity, 2000); ex != 0 {
		return nil, ex
	}
	status, ex := h.ReadDiscreteInputs(ctx, address, quantity)
	switch {
	case ex != 0:
		return nil, ex
	case len(status) != int(quantity):
		return nil, SlaveDeviceFailure
	}
	return put(1+int(byteCount(quantity)), byte(byteCount(quantity)), status), 0
}

func (h *Mux) readHoldingRegisters(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.ReadHoldingRegisters == nil:
		return nil, IllegalFunction
	case len(req) != 4:
		return nil, IllegalDataAddress
	}
	address := binary.BigEndian.Uint16(req[0:])
	quantity := binary.BigEndian.Uint16(req[2:])
	if ex := boundCheck(address, quantity, 125); ex != 0 {
		return nil, ex
	}
	values, ex := h.ReadHoldingRegisters(ctx, address, quantity)
	switch {
	case ex != 0:
		return nil, ex
	case len(values) != 2*int(quantity):
		return nil, SlaveDeviceFailure
	}
	return put(1+int(quantity)*2, byte(quantity*2), values), 0
}

func (h *Mux) readInputRegisters(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.ReadInputRegisters == nil:
		return nil, IllegalFunction
	case len(req) != 4:
		return nil, IllegalDataAddress
	}
	address := binary.BigEndian.Uint16(req[0:])
	quantity := binary.BigEndian.Uint16(req[2:])
	if ex := boundCheck(address, quantity, 125); ex != 0 {
		return nil, ex
	}
	values, ex := h.ReadInputRegisters(ctx, address, quantity)
	switch {
	case ex != 0:
		return nil, ex
	case len(values) != 2*int(quantity):
		return nil, SlaveDeviceFailure
	}
	return put(1+int(quantity)*2, byte(quantity*2), values), 0
}

func (h *Mux) writeSingleCoil(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.WriteSingleCoil == nil:
		return nil, IllegalFunction
	case len(req) != 4:
		return nil, IllegalDataAddress
	}
	address := binary.BigEndian.Uint16(req[0:])
	status := false
	switch binary.BigEndian.Uint16(req[2:]) {
	case 0x0000:
	case 0xFF00:
		status = true
	default:
		return nil, IllegalDataValue
	}
	if ex = h.WriteSingleCoil(ctx, address, status); ex != 0 {
		return nil, ex
	}
	return req, 0
}

func (h *Mux) writeSingleRegister(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.WriteSingleRegister == nil:
		return nil, IllegalFunction
	case len(req) != 4:
		return nil, IllegalDataAddress
	}
	address := binary.BigEndian.Uint16(req[0:])
	value := binary.BigEndian.Uint16(req[2:])
	if ex = h.WriteSingleRegister(ctx, address, value); ex != 0 {
		return nil, ex
	}
	return req, 0
}

func (h *Mux) writeMultipleCoils(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.WriteMultipleCoils == nil:
		return nil, IllegalFunction
	case len(req) < 6:
		return nil, IllegalDataAddress
	}
	address := binary.BigEndian.Uint16(req[0:])
	quantity := binary.BigEndian.Uint16(req[2:])
	if len(req[5:]) != int(req[4]) {
		return nil, IllegalDataValue
	}
	if ex := boundCheck(address, quantity, 1968); ex != 0 {
		return nil, ex
	}
	if ex = h.WriteMultipleCoils(ctx, address, bytesToBools(quantity, req[5:])); ex != 0 {
		return nil, ex
	}
	return req[:4], 0
}

func (h *Mux) writeMultipleRegisters(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.WriteMultipleRegisters == nil:
		return nil, IllegalFunction
	case len(req) < 6:
		return nil, IllegalDataAddress
	}
	address := binary.BigEndian.Uint16(req[0:])
	quantity := binary.BigEndian.Uint16(req[2:])
	if 2*quantity != uint16(req[4]) || int(req[4]) != len(req[5:]) {
		return nil, IllegalDataValue
	}
	if ex := boundCheck(address, quantity, 123); ex != 0 {
		return nil, ex
	}
	if ex = h.WriteMultipleRegisters(ctx, address, req[5:]); ex != 0 {
		return nil, ex
	}
	return req[:4], 0
}

func (h *Mux) readWriteMultipleRegisters(ctx cancel.Context, req []byte) (res []byte, ex Exception) {
	switch {
	case h.ReadWriteMultipleRegisters == nil:
		return nil, IllegalFunction
	case len(req) < 11:
		return nil, IllegalDataAddress
	}
	rAddress := binary.BigEndian.Uint16(req[0:])
	rQuantity := binary.BigEndian.Uint16(req[2:])
	wAddress := binary.BigEndian.Uint16(req[4:])
	wQuantity := binary.BigEndian.Uint16(req[6:])
	if rQuantity*2 != uint16(req[8]) || int(req[8]) != len(req[9:]) {
		return nil, IllegalDataValue
	}
	if ex := boundCheck(rAddress, rQuantity, 125); ex != 0 {
		return nil, ex
	}
	if ex := boundCheck(wAddress, wQuantity, 121); ex != 0 {
		return nil, ex
	}
	res, ex = h.ReadWriteMultipleRegisters(ctx, rAddress, rQuantity, wAddress, req[9:])
	switch {
	case ex != 0:
		return nil, ex
	case len(res) != int(rQuantity)*2:
		return nil, SlaveDeviceFailure
	}
	return put(1+len(res), byte(len(res)), res), 0
}
