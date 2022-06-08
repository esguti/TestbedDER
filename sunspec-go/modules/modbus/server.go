package modbus

import (
	"sync"

	"github.com/GoAethereal/cancel"
)

// Server is the go implementation of a modbus slave.
// Once serving it will listen for incoming requests and forward them to the modbus.Handler h.
// Generally the intended use is as follows:
//	ctx := cancel.New()
//	cfg := modbus.Config{
//		Mode:     "tcp",
//		Kind:     "tcp",
//		Endpoint: "localhost:502",
//	}
//	s := cfg.Server()
//	h := &modbus.Mux{/*define individual handlers*/}
//
//	log.Fatal(s.Serve(ctx,h))
type Server struct {
	cfg Config
	framer
}

// Serve starts the modbus server and listens for incoming requests.
// The Handler h is called for each inbound message.
// h must be safe for use by multiple go routines.
func (s *Server) Serve(ctx cancel.Context, h Handler) error {
	var wg sync.WaitGroup
	l, err := s.cfg.listen(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return nil
		default:
			conn, err := l()
			if err != nil {
				continue
			}
			wg.Add(1)
			go func(conn connection) {
				defer wg.Done()
				s.handle(ctx, conn, h)
			}(conn)
		}
	}
}

// handle starts up a new request handler for a given connection
func (s *Server) handle(ctx cancel.Context, c connection, h Handler) {
	defer c.close()
	var wg sync.WaitGroup

	wait := c.listen(ctx, func(adu []byte, err error) (quit bool) {
		if err != nil {
			return true
		}
		buf := s.buffer()
		buf = buf[:copy(buf, adu)]
		wg.Add(1)
		go func(adu []byte) {
			defer wg.Done()
			var res []byte
			var ex Exception
			code, req, err := s.decode(adu)

			switch {
			case err != nil:
				return
			case code < 0x80:
				res, ex = h.Handle(ctx, code, req)
			default:
				ex = IllegalFunction
			}

			switch {
			case ex != 0:
				code |= 0x80
				res = []byte{byte(ex)}
			case len(res) > 252:
				code |= 0x80
				res = []byte{byte(SlaveDeviceFailure)}
			}

			res, _ = s.reply(code, res, adu)
			if err := c.write(ctx, res); err != nil {
				return
			}
		}(buf)
		return false
	})

	c.read(ctx, s.buffer())
	<-wait
	wg.Wait()
}
