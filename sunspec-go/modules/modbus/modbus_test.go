package modbus_test

import (
	"sync"
	"testing"
	"time"

	"github.com/GoAethereal/cancel"
	"github.com/GoAethereal/modbus"
)

var cfg = modbus.Config{
	Mode:     "tcp",
	Kind:     "tcp",
	Endpoint: "localhost:1337",
}

var (
	mu sync.Mutex
	s  = cfg.Server()
	c  = cfg.Client()
)

func TestReadCoils(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[[2]uint16][]bool{
		{0, 1}:     {true},
		{10, 5}:    {true, false, true, false, true},
		{65535, 1}: {false},
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		ReadCoils: func(ctx cancel.Context, address, quantity uint16) (res []bool, ex modbus.Exception) {
			if res, ok := testCases[[2]uint16{address, quantity}]; ok {
				return res, 0
			}
			return nil, modbus.IllegalDataAddress
		},
	})

	time.Sleep(250 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client connection refused: %v", err)
	}
	defer c.Disconnect()

	for aq, want := range testCases {
		res, err := c.ReadCoils(ctx, aq[0], aq[1])
		if err != nil {
			t.Fatalf("read coils failed: %v", err)
		}
		for i := range res {
			if res[i] != want[i] {
				t.Fatalf("read coils received invalid value at index %v; want %v; got: %v", i, res[i], want[i])
			}
		}
	}
}

func TestReadDiscreteInputs(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[[2]uint16][]bool{
		{0, 1}:     {true},
		{10, 5}:    {true, false, true, false, true},
		{65535, 1}: {false},
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		ReadDiscreteInputs: func(ctx cancel.Context, address, quantity uint16) (res []bool, ex modbus.Exception) {
			if res, ok := testCases[[2]uint16{address, quantity}]; ok {
				return res, 0
			}
			return nil, modbus.IllegalDataAddress
		},
	})

	time.Sleep(250 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client connection refused: %v", err)
	}
	defer c.Disconnect()

	for aq, want := range testCases {
		res, err := c.ReadDiscreteInputs(ctx, aq[0], aq[1])
		if err != nil {
			t.Fatalf("read discrete inputs failed: %v", err)
		}
		for i := range res {
			if res[i] != want[i] {
				t.Fatalf("read discrete inputs received invalid value at index %v; want %v; got: %v", i, res[i], want[i])
			}
		}
	}
}

func TestReadHoldingRegisters(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[[2]uint16][]byte{
		{0, 1}:     {0, 1},
		{10, 5}:    {0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{65535, 1}: {1, 0},
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		ReadHoldingRegisters: func(ctx cancel.Context, address uint16, quantity uint16) (res []byte, ex modbus.Exception) {
			if res, ok := testCases[[2]uint16{address, quantity}]; ok {
				return res, 0
			}
			return nil, modbus.IllegalDataAddress
		},
	})

	time.Sleep(250 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client connection refused: %v", err)
	}
	defer c.Disconnect()

	for aq, want := range testCases {
		res, err := c.ReadHoldingRegisters(ctx, aq[0], aq[1])
		if err != nil {
			t.Fatalf("read holding registers failed: %v", err)
		}
		for i := range res {
			if res[i] != want[i] {
				t.Fatalf("read holding registers received invalid value at index %v; want %v; got: %v", i, res[i], want[i])
			}
		}
	}
}

func TestReadInputRegisters(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[[2]uint16][]byte{
		{0, 1}:     {0, 1},
		{10, 5}:    {0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{65535, 1}: {1, 0},
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		ReadInputRegisters: func(ctx cancel.Context, address uint16, quantity uint16) (res []byte, ex modbus.Exception) {
			if res, ok := testCases[[2]uint16{address, quantity}]; ok {
				return res, 0
			}
			return nil, modbus.IllegalDataAddress
		},
	})

	time.Sleep(250 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client connection refused: %v", err)
	}
	defer c.Disconnect()

	for aq, want := range testCases {
		res, err := c.ReadInputRegisters(ctx, aq[0], aq[1])
		if err != nil {
			t.Fatalf("read input registers failed: %v", err)
		}
		for i := range res {
			if res[i] != want[i] {
				t.Fatalf("read input registers received invalid value at index %v; want %v; got: %v", i, res[i], want[i])
			}
		}
	}
}

func TestWriteSingleCoil(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[uint16]bool{
		0:     true,
		10:    false,
		65535: true,
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		WriteSingleCoil: func(ctx cancel.Context, address uint16, status bool) (ex modbus.Exception) {
			if want, ok := testCases[address]; ok {
				if want != status {
					t.Errorf("server received unexpected value handling function code WriteSingleCoil at address %v, got %v; wanted: %v", address, status, want)
					return modbus.SlaveDeviceFailure
				}
				return 0
			}
			t.Errorf("server received unexpected address %v for handling function code WriteSingleCoil", address)
			return modbus.IllegalDataAddress
		},
	})

	time.Sleep(250 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client connection refused: %v", err)
	}
	defer c.Disconnect()

	for a, want := range testCases {
		if err := c.WriteSingleCoil(ctx, a, want); err != nil {
			t.Fatalf("write single coil failed: %v", err)
		}
	}
}

func TestWriteSingleRegister(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[uint16]uint16{
		0:     666,
		10:    69,
		65535: 42,
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		WriteSingleRegister: func(ctx cancel.Context, address uint16, value uint16) (ex modbus.Exception) {
			if want, ok := testCases[address]; ok {
				if want != value {
					t.Errorf("server received unexpected value handling function code WriteSingleRegister at address %v, got %v; wanted: %v", address, value, want)
					return modbus.SlaveDeviceFailure
				}
				return 0
			}
			t.Errorf("server received unexpected address %v for handling function code WriteSingleRegister", address)
			return modbus.IllegalDataAddress
		},
	})

	time.Sleep(250 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client connection refused: %v", err)
	}
	defer c.Disconnect()

	for a, want := range testCases {
		if err := c.WriteSingleRegister(ctx, a, want); err != nil {
			t.Fatalf("write single register failed: %v", err)
		}
	}
}

func TestWriteMultipleCoils(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[uint16][]bool{
		0:     {true},
		10:    {true, false, true, false, true},
		65535: {false},
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		WriteMultipleCoils: func(ctx cancel.Context, address uint16, status []bool) (ex modbus.Exception) {
			if want, ok := testCases[address]; ok {
				for i := range want {
					if want[i] != status[i] {
						t.Errorf("server received unexpected status handling function code WriteMultipleCoils at address %v, got %v; wanted: %v", address, status, want)
						return modbus.SlaveDeviceFailure
					}
				}
				return 0
			}
			t.Errorf("server received unexpected address %v for handling function code WriteMultipleCoils", address)
			return modbus.IllegalDataAddress
		},
	})

	time.Sleep(1 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client: connection refused: %v", err)
	}
	defer c.Disconnect()

	for a, want := range testCases {
		if err := c.WriteMultipleCoils(ctx, a, want...); err != nil {
			t.Fatalf("client: WriteMultipleCoils failed: %v", err)
		}
	}
}

func TestWriteMultipleRegisters(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[uint16][]byte{
		0:     {0, 1},
		10:    {0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		65535: {0, 1},
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		WriteMultipleRegisters: func(ctx cancel.Context, address uint16, values []byte) (ex modbus.Exception) {
			if want, ok := testCases[address]; ok {
				for i := range want {
					if want[i] != values[i] {
						t.Errorf("server received unexpected status handling function code WriteMultipleRegisters at address %v, got %v; wanted: %v", address, values, want)
						return modbus.SlaveDeviceFailure
					}
				}
				return 0
			}
			t.Errorf("server received unexpected address %v for handling function code WriteMultipleRegisters", address)
			return modbus.IllegalDataAddress
		},
	})

	time.Sleep(1 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client: connection refused: %v", err)
	}
	defer c.Disconnect()

	for a, want := range testCases {
		if err := c.WriteMultipleRegisters(ctx, a, want); err != nil {
			t.Fatalf("client: WriteMultipleRegisters failed: %v", err)
		}
	}
}

func TestReadWriteMultipleRegisters(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	testCases := map[[3]uint16][2][]byte{
		{0, 1, 0}:         {{0, 1}, {1, 0}},
		{10, 5, 10}:       {{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},
		{65535, 1, 65535}: {{1, 0}, {0, 1}},
	}

	ctx := cancel.New()
	defer ctx.Cancel()

	go s.Serve(ctx, &modbus.Mux{
		ReadWriteMultipleRegisters: func(ctx cancel.Context, rAddress uint16, rQuantity uint16, wAddress uint16, values []byte) (res []byte, ex modbus.Exception) {
			if want, ok := testCases[[3]uint16{rAddress, rQuantity, wAddress}]; ok {
				for i, v := range want[1] {
					if v != values[i] {
						t.Errorf("server received unexpected value whilst handling function code ReadWriteMultipleRegisters at address %v, got %v; wanted: %v", wAddress, values, want)
						return nil, modbus.SlaveDeviceFailure
					}
				}
				return want[0], 0
			}
			t.Errorf("server received unexpected request for handling function code ReadWriteMultipleRegisters with read address %v; read quantity %v; write address %v", rAddress, rQuantity, wAddress)
			return nil, modbus.IllegalDataAddress
		},
	})

	time.Sleep(1 * time.Millisecond)

	if err := c.Connect(); err != nil {
		t.Fatalf("client: connection refused: %v", err)
	}
	defer c.Disconnect()

	for dst, want := range testCases {
		res, err := c.ReadWriteMultipleRegisters(ctx, dst[0], dst[1], dst[2], want[1])
		if err != nil {
			t.Fatalf("client: ReadWriteMultipleRegisters failed: %v", err)
		}
		for i, v := range want[0] {
			if v != res[i] {
				t.Fatalf("read input registers received invalid value at index %v; want %v; got: %v", i, res, want[0])
			}
		}
	}
}
