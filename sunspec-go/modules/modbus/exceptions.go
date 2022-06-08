package modbus

import "fmt"

const (
	// IllegalFunction - Exception code 0x01
	//
	// The function code received in the query is not an allowable action for the server (or slave). This
	// may be because the function code is only applicable to newer devices, and was not
	// implemented in the unit selected. It could also indicate that the server (or slave) is in the wrong
	// state to process a request of this type, for example because it is unconfigured and is being
	// asked to return register values.
	IllegalFunction Exception = 0x01
	// IllegalDataAddress - Exception code 0x02
	//
	// The data address received in the query is not an allowable address for the server (or slave). More
	// specifically, the combination of reference number and transfer length is invalid. For a controller with
	// 100 registers, the PDU addresses the first register as 0, and the last one as 99. If a request
	// is submitted with a starting register address of 96 and a quantity of registers of 4, then this request
	// will successfully operate (address-wise at least) on registers 96, 97, 98, 99. If a request is
	// submitted with a starting register address of 96 and a quantity of registers of 5, then this request
	// will fail with Exception Code 0x02 “Illegal Data Address” since it attempts to operate on registers
	// 96, 97, 98, 99 and 100, and there is no register with address 100.
	IllegalDataAddress Exception = 0x02
	// IllegalDataValue - Exception code 0x03
	//
	// A value contained in the query data field is not an allowable value for server (or slave). This
	// indicates a fault in the structure of the remainder of a complex request, such as that the implied
	// length is incorrect. It specifically does NOT mean that a data item submitted for storage in a register
	// has a value outside the expectation of the application program, since the MODBUS protocol
	// is unaware of the significance of any particular value of any particular register.
	IllegalDataValue Exception = 0x03
	// SlaveDeviceFailure - Exception code 0x04
	//
	// An unrecoverable error occurred while the server (or slave) was attempting to perform the
	// requested action.
	SlaveDeviceFailure Exception = 0x04
	// Acknowledge - Exception code 0x05
	//
	// Specialized use in conjunction with programming commands. The server (or slave) has accepted the request
	// and is processing it, but a long duration of time will be required to do so. This response is
	// returned to prevent a timeout error from occurring in the client (or master). The client (or master)
	// can next issue a Poll Program Complete message to determine if processing is completed.
	Acknowledge Exception = 0x05
	// SlaveDeviceBusy - Exception code 0x06
	//
	// Specialized use in conjunction with programming commands. The server (or slave) is engaged in processing a
	// long–duration program command. The client (or master) should retransmit the message later when
	// the server (or slave) is free
	SlaveDeviceBusy Exception = 0x06
	// MemoryParityError - Exception code 0x08
	//
	// Specialized use in conjunction with function codes 20 and 21 and reference type 6, to indicate that
	// the extended file area failed to pass a consistency check. The server (or slave) attempted to read record
	// file, but detected a parity error in the memory. The client (or master) can retry the request, but
	// service may be required on the server (or slave) device.
	MemoryParityError Exception = 0x08
	// GatewayPathUnavailable - Exception code 0x0A
	//
	// Specialized use in conjunction with gateways, indicates that the gateway was unable to allocate
	// an internal communication path from the input port to the output port for processing the request.
	// Usually means that the gateway is misconfigured or overloaded.
	GatewayPathUnavailable Exception = 0x0A
	// GatewayTargetDeviceFailedToRespond - Exception code 0x0B
	//
	// Specialized use in conjunction with gateways, indicates that no response was obtained from the
	// target device. Usually means that the device is not present on the network.
	GatewayTargetDeviceFailedToRespond Exception = 0x0B
)

// Exception represents a modbus exception as defined by the specification.
// It´s a superset of the error interface.
type Exception byte

// Error returns a human readable string representing the underlying exception.
func (ex Exception) Error() string {
	prefix := "modbus: exception "
	switch ex {
	case IllegalFunction:
		return prefix + "illegal function"
	case IllegalDataAddress:
		return prefix + "illegal data address"
	case IllegalDataValue:
		return prefix + "illegal data value"
	case SlaveDeviceFailure:
		return prefix + "slave device failure"
	case Acknowledge:
		return prefix + "acknowledge"
	case SlaveDeviceBusy:
		return prefix + "slave device busy"
	case MemoryParityError:
		return prefix + "memory parity error"
	case GatewayPathUnavailable:
		return prefix + "gateway path unavailable"
	case GatewayTargetDeviceFailedToRespond:
		return prefix + "gateway target device failed to respond"
	}
	return prefix + fmt.Sprintf("code %v undefined", byte(ex))
}
