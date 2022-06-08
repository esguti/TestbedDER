package main

import (
	"fmt"
	"time"

	"github.com/GoAethereal/cancel"
	"github.com/GoAethereal/modbus"
)

var cfg = modbus.Config{
	Mode:     "tcp",
	Kind:     "tcp",
	Endpoint: "localhost:1337",
}

var ctx = cancel.New()

func main() {
	// start the modbus master
	go client()
	// start the modbus slave
	server()
}

func server() {
	// instantiate a new modbus slave from the given configuration
	s := cfg.Server()
	// start serving the function codes read coils and write multiple registers
	s.Serve(ctx, &modbus.Mux{
		// ReadCoils will always respond with an alternating pattern of coil states
		ReadCoils: func(ctx cancel.Context, address, quantity uint16) (res []bool, ex modbus.Exception) {
			res = make([]bool, quantity)
			for i := range res {
				res[i] = i%2 == 0
			}
			return res, 0
		},
		// WriteMultipleRegisters will print out the received register values as string
		WriteMultipleRegisters: func(ctx cancel.Context, address uint16, values []byte) (ex modbus.Exception) {
			fmt.Printf("server: received write multiple registers request: %v\n", string(values))
			return 0
		},
	})
}

func client() {
	// cancel the context when the client is done -> this will also initialize a server shutdown
	defer ctx.Cancel()
	// instantiate a new modbus master from the given configuration
	c := cfg.Client()
	// give the server some time to start up
	time.Sleep(1 * time.Second)
	// attempt to connect to the server
	if err := c.Connect(); err != nil {
		return
	}
	defer c.Disconnect()
	for i := 0; i < 10; i++ {
		// request 10 coil states from the server
		res, err := c.ReadCoils(ctx, 0, 10)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("client: received response for read coils %v\n", res)
		time.Sleep(1 * time.Second)
		// write hello there to the server
		if err := c.WriteMultipleRegisters(ctx, 0, []byte("hello there!")); err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}
