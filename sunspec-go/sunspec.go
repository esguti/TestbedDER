package sunspec

// Logger defines the behavior for the internally used, optionally provided logger.
type Logger interface {
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
}

type logger struct{}

// Debug voids all args.
func (l logger) Debug(_ ...interface{}) {}

// Error voids all args.
func (l logger) Error(_ ...interface{}) {}

// Warn voids all args.
func (l logger) Warn(_ ...interface{}) {}

// Info voids all args.
func (l logger) Info(_ ...interface{}) {}

// marker returns a dummy model for representing the magic identifier SunS.
func marker(adr uint16) Model {
	return &model{
		&group{
			name: "marker",
			points: Points{
				&tString{
					data: []byte("SunS"),
					point: point{
						name:    "SunS",
						static:  true,
						address: adr,
					},
				},
			},
		},
	}
}

// header returns a prototype for identifying a model using the minimum requirements.
func header(adr, id, l uint16) Model {
	return &model{
		&group{
			name: "header",
			points: Points{
				&tUint16{
					data: id,
					point: point{
						name:    "ID",
						static:  true,
						address: adr,
					},
				},
				&tUint16{
					data: l,
					point: point{
						name:    "L",
						static:  true,
						address: adr + 1,
					},
				},
			},
		},
	}
}
