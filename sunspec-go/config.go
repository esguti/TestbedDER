package sunspec

// Config is the configuration for a client or server.
type Config struct {
	// Endpoint specifics the sunspec host and is mandatory.
	// Currently only modbus tcp-networking is supported.
	// The schema must be host:port
	Endpoint string
	// Logger can be optionally defined.
	Logger Logger
}

// logger returns the optional logger.
// If none is defined (nil) the default logger is returned.
func (o *Config) logger() Logger {
	if o.Logger != nil {
		return o.Logger
	}
	return logger{}
}

// Client instantiates a new client from the given configuration.
func (o Config) Client() *Client {
	return &Client{client: newModbusClient(o.Endpoint, o.logger()), logger: o.logger()}
}

// Server instantiates a new server from the given configuration.
func (o Config) Server() *Server {
	return &Server{server: newModbusServer(o.Endpoint, o.logger()), logger: o.logger()}
}
