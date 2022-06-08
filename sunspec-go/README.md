# Sunspec

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/TRICERA-energy/sunspec)](https://goreportcard.com/report/github.com/TRICERA-energy/sunspec)
[![Go Reference](https://pkg.go.dev/badge/github.com/TRICERA-energy/sunspec.svg)](https://pkg.go.dev/github.com/TRICERA-energy/sunspec)

This package provides functions and definitions for creating [sunspec](https://sunspec.org/) compliant server and client devices using the go programming language.

Examples for using the package can be found [here](https://github.com/TRICERA-energy/sunspec/tree/master/examples).

**NOTICE: Currently only communication via modbus-TCP is supported.**

## Type system

Data types defined by the sunspec specification are represented in this library using their own custom interface. The package guarantees that point-type interfaces provided by the client or server also satisfy one of the type interfaces. This way assertion can be used to get explicit access to specific functionalities. 

The following types are available:

```go
sunspec.Int16
sunspec.Int32
sunspec.Int64
sunspec.Pad
sunspec.Sunssf
sunspec.Uint16
sunspec.Uint32
sunspec.Uint64
sunspec.Acc16
sunspec.Acc32
sunspec.Acc64
sunspec.Count
sunspec.Bitfield16
sunspec.Bitfield32
sunspec.Bitfield64
sunspec.Enum16
sunspec.Enum32
sunspec.String
sunspec.Float32
sunspec.Float64
sunspec.Ipaddr
sunspec.Ipv6addr
sunspec.Eui48
```