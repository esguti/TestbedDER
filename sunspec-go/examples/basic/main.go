package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/GoAethereal/cancel"
	"github.com/TRICERA-energy/sunspec"
)

var port = flag.Int("Port", 1337, "Port the sunspec communication should use")

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	ctx    = cancel.New()
)

var (
	defs     []sunspec.Definition
	endpoint string
)

func main() {
	flag.Parse()

	endpoint = fmt.Sprintf("192.168.1.13:%v", *port)

	for _, model := range models {
		var def sunspec.ModelDef
		// unmarshal a model schema into it´s go definition
		if err := json.Unmarshal(model, &def); err != nil {
			logger.Fatalln(err)
		}
		defs = append(defs, &def)
	}

	var wg sync.WaitGroup

	// start the server
	wg.Add(1)
	go func() {
		defer wg.Done()
		Server()
	}()

	time.Sleep(1 * time.Second)

	//start the client
	wg.Add(1)
	go func() {
		defer wg.Done()
		Client()
	}()

	wg.Wait()
}

// handler gets called for any incoming sunspec request
func handler(ctx cancel.Context, req sunspec.Request) error {
	defer req.Flush()
	for _, p := range req.Points() {
		if p, ok := p.(sunspec.Float32); ok {
			p.Set(rand.Float32())
		}
	}
	return nil
}

// Server starts up a new sunspec server.
func Server() {
	// create a new sunspec server instance
	s := (sunspec.Config{Endpoint: endpoint}).Server()

	// start serving
	logger.Println(s.Serve(ctx, handler, defs...))
}

// Client instantiates a new sunspec client.
// First scanning then polling the server every second.
func Client() {
	// create a new client requesting data from the server
	c := (sunspec.Config{Endpoint: endpoint}).Client()

	// attempt to connect to the server
	if err := c.Connect(); err != nil {
		logger.Fatalln(err)
	}
	defer c.Disconnect()

	// scan the endpoint retrieving all models
	if err := c.Scan(ctx, defs...); err != nil {
		logger.Fatalln("Read error:", err)
	}

	// continuously read the entire model from the server
	for {
		for _, m := range c.Models() {
			if _, err := c.Read(ctx, m); err != nil {
				logger.Println(err)
			}
			fmt.Printf("Read Model %v:\n", m.ID())
			for _, p := range m.Points() {
				fmt.Printf("\t%v\t- current value: %v\n", p.Name(), p)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

var models = [][]byte{
	[]byte(`{
    "group": {
        "desc": "All SunSpec compliant devices must include this as the first model",
        "label": "Common",
        "name": "common",
        "points": [
            {
                "desc": "Model identifier",
                "label": "Model ID",
                "mandatory": "M",
                "name": "ID",
                "size": 1,
                "static": "S",
                "type": "uint16",
                "value": 1
            },
            {
                "desc": "Model length",
                "label": "Model Length",
                "mandatory": "M",
                "name": "L",
                "size": 1,
                "static": "S",
                "type": "uint16",
                "value": 66
            },
            {
                "desc": "Well known value registered with SunSpec for compliance",
                "label": "Manufacturer",
                "mandatory": "M",
                "name": "Mn",
                "size": 16,
                "static": "S",
                "type": "string",
                "value": "TRICERA energy"
            },
            {
                "desc": "Manufacturer specific value (32 chars)",
                "label": "Model",
                "mandatory": "M",
                "name": "Md",
                "size": 16,
                "static": "S",
                "type": "string",
                "value": "example"
            },
            {
                "desc": "Manufacturer specific value (16 chars)",
                "label": "Options",
                "name": "Opt",
                "size": 8,
                "static": "S",
                "type": "string",
                "value": "none"
            },
            {
                "desc": "Manufacturer specific value (16 chars)",
                "label": "Version",
                "name": "Vr",
                "size": 8,
                "static": "S",
                "type": "string",
                "value": "v0.0.1"
            },
            {
                "desc": "Manufacturer specific value (32 chars)",
                "label": "Serial Number",
                "mandatory": "M",
                "name": "SN",
                "size": 16,
                "static": "S",
                "type": "string"
            },
            {
                "access": "RW",
                "desc": "Modbus device address",
                "label": "Device Address",
                "name": "DA",
                "size": 1,
                "type": "uint16"
            },
            {
                "desc": "Force even alignment",
                "name": "Pad",
                "size": 1,
                "static": "S",
                "type": "pad"
            }
        ],
        "type": "group"
    },
    "id": 1
}`),
	[]byte(`{
    "group": {
        "desc": "Include this model for three phase inverter monitoring using float values",
        "label": "Inverter (Three Phase) FLOAT",
        "name": "inverter",
        "points": [
            {
                "desc": "Model identifier",
                "label": "Model ID",
                "mandatory": "M",
                "name": "ID",
                "size": 1,
                "static": "S",
                "type": "uint16",
                "value": 113
            },
            {
                "desc": "Model length",
                "label": "Model Length",
                "mandatory": "M",
                "name": "L",
                "size": 1,
                "static": "S",
                "type": "uint16",
                "value": 60
            },
            {
                "desc": "AC Current",
                "label": "Amps",
                "mandatory": "M",
                "name": "A",
                "size": 2,
                "type": "float32",
                "units": "A"
            },
            {
                "desc": "Phase A Current",
                "label": "Amps PhaseA",
                "mandatory": "M",
                "name": "AphA",
                "size": 2,
                "type": "float32",
                "units": "A"
            },
            {
                "desc": "Phase B Current",
                "label": "Amps PhaseB",
                "mandatory": "M",
                "name": "AphB",
                "size": 2,
                "type": "float32",
                "units": "A"
            },
            {
                "desc": "Phase C Current",
                "label": "Amps PhaseC",
                "mandatory": "M",
                "name": "AphC",
                "size": 2,
                "type": "float32",
                "units": "A"
            },
            {
                "desc": "Phase Voltage AB",
                "label": "Phase Voltage AB",
                "name": "PPVphAB",
                "size": 2,
                "type": "float32",
                "units": "V"
            },
            {
                "desc": "Phase Voltage BC",
                "label": "Phase Voltage BC",
                "name": "PPVphBC",
                "size": 2,
                "type": "float32",
                "units": "V"
            },
            {
                "desc": "Phase Voltage CA",
                "label": "Phase Voltage CA",
                "name": "PPVphCA",
                "size": 2,
                "type": "float32",
                "units": "V"
            },
            {
                "desc": "Phase Voltage AN",
                "label": "Phase Voltage AN",
                "mandatory": "M",
                "name": "PhVphA",
                "size": 2,
                "type": "float32",
                "units": "V"
            },
            {
                "desc": "Phase Voltage BN",
                "label": "Phase Voltage BN",
                "mandatory": "M",
                "name": "PhVphB",
                "size": 2,
                "type": "float32",
                "units": "V"
            },
            {
                "desc": "Phase Voltage CN",
                "label": "Phase Voltage CN",
                "mandatory": "M",
                "name": "PhVphC",
                "size": 2,
                "type": "float32",
                "units": "V"
            },
            {
                "desc": "AC Power",
                "label": "Watts",
                "mandatory": "M",
                "name": "W",
                "size": 2,
                "type": "float32",
                "units": "W"
            },
            {
                "desc": "Line Frequency",
                "label": "Hz",
                "mandatory": "M",
                "name": "Hz",
                "size": 2,
                "type": "float32",
                "units": "Hz"
            },
            {
                "desc": "AC Apparent Power",
                "label": "VA",
                "name": "VA",
                "size": 2,
                "type": "float32",
                "units": "VA"
            },
            {
                "desc": "AC Reactive Power",
                "label": "VAr",
                "name": "VAr",
                "size": 2,
                "type": "float32",
                "units": "var"
            },
            {
                "desc": "AC Power Factor",
                "label": "PF",
                "name": "PF",
                "size": 2,
                "type": "float32",
                "units": "Pct"
            },
            {
                "desc": "AC Energy",
                "label": "WattHours",
                "mandatory": "M",
                "name": "WH",
                "size": 2,
                "type": "float32",
                "units": "Wh"
            },
            {
                "desc": "DC Current",
                "label": "DC Amps",
                "name": "DCA",
                "size": 2,
                "type": "float32",
                "units": "A"
            },
            {
                "desc": "DC Voltage",
                "label": "DC Voltage",
                "name": "DCV",
                "size": 2,
                "type": "float32",
                "units": "V"
            },
            {
                "desc": "DC Power",
                "label": "DC Watts",
                "name": "DCW",
                "size": 2,
                "type": "float32",
                "units": "W"
            },
            {
                "desc": "Cabinet Temperature",
                "label": "Cabinet Temperature",
                "mandatory": "M",
                "name": "TmpCab",
                "size": 2,
                "type": "float32",
                "units": "C"
            },
            {
                "desc": "Heat Sink Temperature",
                "label": "Heat Sink Temperature",
                "name": "TmpSnk",
                "size": 2,
                "type": "float32",
                "units": "C"
            },
            {
                "desc": "Transformer Temperature",
                "label": "Transformer Temperature",
                "name": "TmpTrns",
                "size": 2,
                "type": "float32",
                "units": "C"
            },
            {
                "desc": "Other Temperature",
                "label": "Other Temperature",
                "name": "TmpOt",
                "size": 2,
                "type": "float32",
                "units": "C"
            },
            {
                "desc": "Enumerated value.  Operating state",
                "label": "Operating State",
                "mandatory": "M",
                "name": "St",
                "size": 1,
                "symbols": [
                    {
                        "name": "OFF",
                        "value": 1
                    },
                    {
                        "name": "SLEEPING",
                        "value": 2
                    },
                    {
                        "name": "STARTING",
                        "value": 3
                    },
                    {
                        "name": "MPPT",
                        "value": 4
                    },
                    {
                        "name": "THROTTLED",
                        "value": 5
                    },
                    {
                        "name": "SHUTTING_DOWN",
                        "value": 6
                    },
                    {
                        "name": "FAULT",
                        "value": 7
                    },
                    {
                        "name": "STANDBY",
                        "value": 8
                    }
                ],
                "type": "enum16"
            },
            {
                "desc": "Vendor specific operating state code",
                "label": "Vendor Operating State",
                "name": "StVnd",
                "size": 1,
                "type": "enum16"
            },
            {
                "desc": "Bitmask value. Event fields",
                "label": "Event1",
                "mandatory": "M",
                "name": "Evt1",
                "size": 2,
                "symbols": [
                    {
                        "name": "GROUND_FAULT",
                        "value": 0
                    },
                    {
                        "name": "DC_OVER_VOLT",
                        "value": 1
                    },
                    {
                        "name": "AC_DISCONNECT",
                        "value": 2
                    },
                    {
                        "name": "DC_DISCONNECT",
                        "value": 3
                    },
                    {
                        "name": "GRID_DISCONNECT",
                        "value": 4
                    },
                    {
                        "name": "CABINET_OPEN",
                        "value": 5
                    },
                    {
                        "name": "MANUAL_SHUTDOWN",
                        "value": 6
                    },
                    {
                        "name": "OVER_TEMP",
                        "value": 7
                    },
                    {
                        "name": "OVER_FREQUENCY",
                        "value": 8
                    },
                    {
                        "name": "UNDER_FREQUENCY",
                        "value": 9
                    },
                    {
                        "name": "AC_OVER_VOLT",
                        "value": 10
                    },
                    {
                        "name": "AC_UNDER_VOLT",
                        "value": 11
                    },
                    {
                        "name": "BLOWN_STRING_FUSE",
                        "value": 12
                    },
                    {
                        "name": "UNDER_TEMP",
                        "value": 13
                    },
                    {
                        "name": "MEMORY_LOSS",
                        "value": 14
                    },
                    {
                        "name": "HW_TEST_FAILURE",
                        "value": 15
                    }
                ],
                "type": "bitfield32"
            },
            {
                "desc": "Reserved for future use",
                "label": "Event Bitfield 2",
                "mandatory": "M",
                "name": "Evt2",
                "size": 2,
                "type": "bitfield32"
            },
            {
                "desc": "Vendor defined events",
                "label": "Vendor Event Bitfield 1",
                "name": "EvtVnd1",
                "size": 2,
                "type": "bitfield32"
            },
            {
                "desc": "Vendor defined events",
                "label": "Vendor Event Bitfield 2",
                "name": "EvtVnd2",
                "size": 2,
                "type": "bitfield32"
            },
            {
                "desc": "Vendor defined events",
                "label": "Vendor Event Bitfield 3",
                "name": "EvtVnd3",
                "size": 2,
                "type": "bitfield32"
            },
            {
                "desc": "Vendor defined events",
                "label": "Vendor Event Bitfield 4",
                "name": "EvtVnd4",
                "size": 2,
                "type": "bitfield32"
            }
        ],
        "type": "group"
    },
    "id": 113
}`)}
