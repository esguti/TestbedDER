# Modbus

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/GoAethereal/modbus)](https://goreportcard.com/report/github.com/GoAethereal/modbus)
[![Go Reference](https://pkg.go.dev/badge/github.com/GoAethereal/modbus.svg)](https://pkg.go.dev/github.com/GoAethereal/modbus)

This package provides functions and definitions for communicating with modbus compliant devices. The implementation tries to follow the [specification](https://modbus.org/docs/Modbus_Application_Protocol_V1_1b.pdf) as closely as possible. However certain parts of the terminology have been changed to feel more familiar to go programmers.  
The focus lies on offering an easy to use client and server.  
  
**NOTICE:**  

The package is currently in alpha state. Significant changes to the signatures are likely in the future.

## Supported features

The following points are available for the server **and** client implementation.  

* context support 
* TCP networking
* modbus TCP payload framing
* asynchronous communication in TCP-framing mode
* function code 0x01: Read Coils
* function code 0x02: Read Discrete Inputs
* function code 0x03: Read Holding Registers
* function code 0x04: Read Input Registers
* function code 0x05: Write Single Coil
* function code 0x06: Write Single Register
* function code 0x0F: Write Multiple Coils
* function code 0x10: Write Multiple Registers
* function code 0x17: Read/Write Multiple Registers

These functionalities are yet to be implemented: 

* serial networking
* UDP networking
* modbus RTU payload framing
* modbus ASCII payload framing
* function code 0x07: Read Exception Status
* function code 0x08: Diagnostics
* function code 0x0B: Get Comm Event Counter
* function code 0x0C: Get Comm Event Log
* function code 0x11: Report Slave ID
* function code 0x14: Read File Record
* function code 0x15: Write File Record
* function code 0x16: Mask Write Registers
* function code 0x18: Read FIFO Queue
* function code 0x2B: Encapsulated Interface Transport

## Installation

Use the following command in a go mod initialized project.

```
go get github.com/GoAethereal/modbus
```
