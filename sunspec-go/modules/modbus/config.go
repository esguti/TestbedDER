package modbus

import (
	"net"

	"github.com/GoAethereal/cancel"
)

// Config are used to configure a modbus client or server
type Config struct {
	// Mode defines the communication framing
	// valid modes are:
	//	- tcp
	//	- rtu	(ToDo)
	//	- ascii	(ToDo)
	Mode string
	// Kind specifies the underlying network layer
	// valid kinds are:
	//	- tcp
	//	- udp 		(ToDo)
	//	- serial	(ToDo)
	Kind string
	// Endpoint used for connecting to (client) or listening on (server)
	Endpoint string
}

// Verify validates the modbus.Options, thereby checking for invalid parameter.
// If the options are valid no error (nil) is returned.
func (cfg *Config) Verify() error {
	switch cfg.Mode {
	case "tcp" /*, "rtu", "ascii"*/ :
	default:
		return ErrInvalidParameter
	}

	switch cfg.Kind {
	case "tcp" /*, "udp", "serial"*/ :
	default:
		return ErrInvalidParameter
	}

	return nil
}

// framer creates a new modbus framer from the given configuration.
func (cfg Config) framer() framer {
	switch cfg.Mode {
	case "tcp":
		return &tcp{}
	}
	return nil
}

// Client instantiates a new modbus master instance from the given configuration.
// If the configuration is malformed nil is returned instead.
// To check the validity of the config use config.Verify()
func (cfg Config) Client() *Client {
	if err := cfg.Verify(); err != nil {
		return nil
	}
	return &Client{cfg: cfg, framer: cfg.framer()}
}

// Server instantiates a new modbus slave instance from the given configuration.
// If the configuration is malformed nil is returned instead.
// To check the validity of the config use config.Verify()
func (cfg Config) Server() *Server {
	if err := cfg.Verify(); err != nil {
		return nil
	}
	return &Server{cfg: cfg, framer: cfg.framer()}
}

// dial attempts to dial in the configured endpoint.
// On success it will return the connection, otherwise an error.
func (cfg Config) dial() (connection, error) {
	switch cfg.Kind {
	case "tcp":
		conn, err := (&net.Dialer{}).Dial(cfg.Kind, cfg.Endpoint)
		if err != nil {
			return nil, err
		}
		return &network{conn: conn}, nil
	}
	return nil, ErrInvalidParameter
}

// listen creates a new listener on the configured endpoint.
// If successfull a acceptor function will be returned.
// The function will block until a new connection is established or an error occurs.
func (cfg Config) listen(ctx cancel.Context) (fn func() (connection, error), err error) {
	switch cfg.Kind {
	case "tcp":
		l, err := net.Listen(cfg.Kind, cfg.Endpoint)
		if err != nil {
			return nil, err
		}
		// start the watch-dog which will stop the listener when the context is canceled
		go func() {
			<-ctx.Done()
			l.Close()
		}()
		fn = func() (connection, error) {
			conn, err := l.Accept()
			return &network{conn: conn}, err
		}

	}
	return fn, nil
}
