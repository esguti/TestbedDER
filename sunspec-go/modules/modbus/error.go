package modbus

import "errors"

var (
	// ErrMissmatchedTransactionId indicates that a received modbus server response did not match
	// the expected request. This error can only occure in TCP-framing mode, where it´s allowed
	// to send multiple parallel request, without waiting for each response.
	// The individual requests are identified by their transaction id.
	// A missmatch therefore means that the response was not intended for the message in question.
	// As a result it should be waited for the next response instead.
	//	NOTICE:
	// This error is handled inside the modbus package and will never be escalated outside the
	// package scope.
	ErrMissmatchedTransactionId = errors.New("modbus: missmatch of transaction id")
	// ErrMissmatchedProtocolId signals a missmatch of the protocol identifier field.
	// A normal response is expected to this value copied from the request.
	ErrMissmatchedProtocolId = errors.New("modbus: missmatch of protocol id")
	// ErrMissmatchedUnitId signals a missmatch of the unit identifier field.
	// A normal response is expected to this value copied from the request.
	ErrMissmatchedUnitId = errors.New("modbus: missmatch of unit id")
	// ErrDataSizeExceeded indicates that the given data length exceeds the limits of a modbus
	// package paylod.
	ErrDataSizeExceeded = errors.New("modbus: data size exceeds limit")
	// ErrInvalidParameter signals a malformed input.
	ErrInvalidParameter = errors.New("modbus: given parameter violates restriction")
)
