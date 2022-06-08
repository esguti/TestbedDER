package sunspec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
)

// Scalable defines the behavior of a point type which may be scaled using the definition:
//	ScaledValue = PointValue * (10^ScaleFactor)
type Scalable interface {
	Scaled() bool
	Factor() int16
}

// scale is internally used to store a scale factor
type scale struct {
	f interface{}
}

// Scaled specifies whether the point is scaled using an optional factor.
func (s *scale) Scaled() bool {
	return s.f != nil
}

// factor returns the scale value of the point.
func (s *scale) factor(p Point) int16 {
	switch sf := s.f.(type) {
	case int16:
		return sf
	case Sunssf:
		return sf.Get()
	case string:
		for g := p.Origin(); g != nil; g = g.Origin() {
			for _, p := range g.Points() {
				if p.Name() == sf {
					if p, ok := p.(Sunssf); ok {
						s.f = p
						return p.Get()
					}
				}
			}
		}
	}
	return 0
}

// ****************************************************************************

// Int16 represents the sunspec type int16.
type Int16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v int16) error
	// Get returns the point´s underlying value.
	Get() int16
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type tInt16 struct {
	point
	data int16
	scale
}

var _ Int16 = (*tInt16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tInt16) Valid() bool { return t.Get() != -0x8000 }

// String formats the point´s value as string.
func (t *tInt16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tInt16) Quantity() uint16 { return 1 }

// encode puts the point´s value into a buffer.
func (t *tInt16) encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, uint16(t.Get()))
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tInt16) decode(buf []byte) error {
	return t.Set(int16(binary.BigEndian.Uint16(buf)))
}

// Set sets the point´s underlying value.
func (t *tInt16) Set(v int16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tInt16) Get() int16 { return t.data }

// Factor returns the scale value of the point.
func (t *tInt16) Factor() int16 { return t.factor(t) }

// Value returns the scaled value as defined by the specification.
func (t *tInt16) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Int32 represents the sunspec type int32.
type Int32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v int32) error
	// Get returns the point´s underlying value.
	Get() int32
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type tInt32 struct {
	point
	data int32
	scale
}

var _ (Int32) = (*tInt32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tInt32) Valid() bool { return t.Get() != -0x80000000 }

// String formats the point´s value as string.
func (t *tInt32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tInt32) Quantity() uint16 { return 2 }

// encode puts the point´s value into a buffer.
func (t *tInt32) encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, uint32(t.Get()))
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tInt32) decode(buf []byte) error {
	return t.Set(int32((binary.BigEndian.Uint32(buf))))
}

// Set sets the point´s underlying value.
func (t *tInt32) Set(v int32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tInt32) Get() int32 { return t.data }

// Factor returns the scale value of the point.
func (t *tInt32) Factor() int16 { return t.factor(t) }

// Value returns the scaled value as defined by the specification.
func (t *tInt32) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Int64 represents the sunspec type int64.
type Int64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v int64) error
	// Get returns the point´s underlying value.
	Get() int64
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type tInt64 struct {
	point
	data int64
	scale
}

var _ Int64 = (*tInt64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tInt64) Valid() bool { return t.Get() != -0x8000000000000000 }

// String formats the point´s value as string.
func (t *tInt64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tInt64) Quantity() uint16 { return 4 }

// encode puts the point´s value into a buffer.
func (t *tInt64) encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, uint64(t.Get()))
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tInt64) decode(buf []byte) error {
	return t.Set(int64(binary.BigEndian.Uint64(buf)))
}

// Set sets the point´s underlying value.
func (t *tInt64) Set(v int64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tInt64) Get() int64 { return t.data }

// Factor returns the scale value of the point.
func (t *tInt64) Factor() int16 { return t.factor(t) }

// Value returns the scaled value as defined by the specification.
func (t *tInt64) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Pad represents the sunspec type pad.
type Pad interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
}

type tPad struct {
	point
}

var _ Pad = (*tPad)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tPad) Valid() bool { return false }

// String formats the point´s value as string.
func (t *tPad) String() string { return "" }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tPad) Quantity() uint16 { return 1 }

// encode puts the point´s value into a buffer.
func (t *tPad) encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, 0x8000)
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tPad) decode(buf []byte) error { return nil }

// ****************************************************************************

// Sunssf represents the sunspec type sunssf.
type Sunssf interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Get returns the point´s underlying value.
	Get() int16
}

type tSunssf struct {
	point
	data int16
}

var _ Sunssf = (*tSunssf)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tSunssf) Valid() bool { return t.Get() != -0x8000 }

// String formats the point´s value as string.
func (t *tSunssf) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tSunssf) Quantity() uint16 { return 1 }

// encode puts the point´s value into a buffer.
func (t *tSunssf) encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, uint16(t.Get()))
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tSunssf) decode(buf []byte) error {
	return t.set(int16(binary.BigEndian.Uint16(buf)))
}

// set sets the point´s underlying value.
func (t *tSunssf) set(v int16) error {
	if v < -10 || v > 10 {
		return errors.New("sunspec: value out of boundary")
	}
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tSunssf) Get() int16 { return t.data }

// ****************************************************************************

// Uint16 represents the sunspec type uint16.
type Uint16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint16) error
	// Get returns the point´s underlying value.
	Get() uint16
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type tUint16 struct {
	point
	data uint16
	scale
}

var _ Uint16 = (*tUint16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tUint16) Valid() bool { return t.Get() != 0xFFFF }

// String formats the point´s value as string.
func (t *tUint16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tUint16) Quantity() uint16 { return 1 }

// encode puts the point´s value into a buffer.
func (t *tUint16) encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tUint16) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *tUint16) Set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tUint16) Get() uint16 { return t.data }

// Factor returns the scale value of the point.
func (t *tUint16) Factor() int16 { return t.factor(t) }

// Value returns the scaled value as defined by the specification.
func (t *tUint16) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Uint32 represents the sunspec type uint32.
type Uint32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint32) error
	// Get returns the point´s underlying value.
	Get() uint32
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type tUint32 struct {
	point
	data uint32
	scale
}

var _ Uint32 = (*tUint32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tUint32) Valid() bool { return t.Get() != 0xFFFFFFFF }

// String formats the point´s value as string.
func (t *tUint32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tUint32) Quantity() uint16 { return 2 }

// encode puts the point´s value into a buffer.
func (t *tUint32) encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tUint32) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint32(buf))
}

// Set sets the point´s underlying value.
func (t *tUint32) Set(v uint32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tUint32) Get() uint32 { return t.data }

// Factor returns the scale value of the point.
func (t *tUint32) Factor() int16 { return t.factor(t) }

// Value returns the scaled value as defined by the specification.
func (t *tUint32) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Uint64 represents the sunspec type uint64.
type Uint64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint64) error
	// Get returns the point´s underlying value.
	Get() uint64
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type tUint64 struct {
	point
	data uint64
	scale
}

var _ Uint64 = (*tUint64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tUint64) Valid() bool { return t.Get() != 0xFFFFFFFFFFFFFFFF }

// String formats the point´s value as string.
func (t *tUint64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tUint64) Quantity() uint16 { return 4 }

// encode puts the point´s value into a buffer.
func (t *tUint64) encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tUint64) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint64(buf))
}

// Set sets the point´s underlying value.
func (t *tUint64) Set(v uint64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tUint64) Get() uint64 { return t.data }

// Factor returns the scale value of the point.
func (t *tUint64) Factor() int16 { return t.factor(t) }

// Value returns the scaled value as defined by the specification.
func (t *tUint64) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Acc16 represents the sunspec type acc16.
type Acc16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint16) error
	// Get returns the point´s underlying value.
	Get() uint16
}

type tAcc16 struct {
	point
	data uint16
	scale
}

var _ Acc16 = (*tAcc16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tAcc16) Valid() bool { return t.Get() != 0 }

// String formats the point´s value as string.
func (t *tAcc16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tAcc16) Quantity() uint16 { return 1 }

// encode puts the point´s value into a buffer.
func (t *tAcc16) encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tAcc16) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *tAcc16) Set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tAcc16) Get() uint16 { return t.data }

// Factor returns the scale value of the point.
func (t *tAcc16) Factor() int16 { return t.factor(t) }

// ****************************************************************************

// Acc32 represents the sunspec type acc32.
type Acc32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint32) error
	// Get returns the point´s underlying value.
	Get() uint32
}

type tAcc32 struct {
	point
	data uint32
	scale
}

var _ (Acc32) = (*tAcc32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tAcc32) Valid() bool { return t.Get() != 0 }

// String formats the point´s value as string.
func (t *tAcc32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tAcc32) Quantity() uint16 { return 2 }

// encode puts the point´s value into a buffer.
func (t *tAcc32) encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tAcc32) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint32(buf))
}

// Set sets the point´s underlying value.
func (t *tAcc32) Set(v uint32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tAcc32) Get() uint32 { return t.data }

// Factor returns the scale value of the point.
func (t *tAcc32) Factor() int16 { return t.factor(t) }

// ****************************************************************************

// Acc64 represents the sunspec type acc64.
type Acc64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint64) error
	// Get returns the point´s underlying value.
	Get() uint64
}

type tAcc64 struct {
	point
	data uint64
	scale
}

var _ Acc64 = (*tAcc64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tAcc64) Valid() bool { return t.Get() != 0 }

// String formats the point´s value as string.
func (t *tAcc64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tAcc64) Quantity() uint16 { return 4 }

// encode puts the point´s value into a buffer.
func (t *tAcc64) encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tAcc64) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint64(buf))
}

// Set sets the point´s underlying value.
func (t *tAcc64) Set(v uint64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tAcc64) Get() uint64 { return t.data }

// Factor returns the scale value of the point.
func (t *tAcc64) Factor() int16 { return t.factor(t) }

// ****************************************************************************

// Count represents the sunspec type count.
type Count interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Get returns the point´s underlying value.
	Get() uint16
}

type tCount struct {
	point
	data uint16
}

var _ Count = (*tCount)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tCount) Valid() bool { return t.Get() != 0 }

// String formats the point´s value as string.
func (t *tCount) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tCount) Quantity() uint16 { return 1 }

// encode puts the point´s value into a buffer.
func (t *tCount) encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tCount) decode(buf []byte) error {
	return t.set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *tCount) set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tCount) Get() uint16 { return t.data }

// ****************************************************************************

// Bitfield16 represents the sunspec type bitfield16.
type Bitfield16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint16) error
	// Get returns the point´s underlying value.
	Get() uint16
	// Flip sets the bit at position pos, starting at 0, to the value of v.
	Flip(pos int, v bool) error
	// Field returns the individual bit values as bool array.
	Field() [16]bool
	// States returns all active enumerated states, correlating the bit value to its symbol.
	States() []string
}

type tBitfield16 struct {
	point
	data    uint16
	symbols Symbols
}

var _ Bitfield16 = (*tBitfield16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tBitfield16) Valid() bool { return t.Get() != 0xFFFF }

// String formats the point´s value as string.
func (t *tBitfield16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tBitfield16) Quantity() uint16 { return 1 }

// encode puts the point´s value into a buffer.
func (t *tBitfield16) encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tBitfield16) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *tBitfield16) Set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tBitfield16) Get() uint16 { return t.data }

// Flip sets the bit at position pos, starting at 0, to the value of v.
func (t *tBitfield16) Flip(pos int, v bool) error {
	switch {
	case pos < 0 || pos > 15:
		return errors.New("sunspec: out of bounds bit position ")
	case v:
		return t.Set(t.Get() | (1 << pos))
	}
	return t.Set(t.Get() &^ (1 << pos))
}

// Field returns the individual bit values as bool array.
func (t *tBitfield16) Field() (f [16]bool) {
	for v, b := t.Get(), 0; b < len(f); b++ {
		f[b] = v&(1<<b) != 0
	}
	return f
}

// States returns all active enumerated states, correlating the bit value to its symbol.
func (t *tBitfield16) States() (s []string) {
	if !t.Valid() {
		return nil
	}
	for i, v := range t.Field() {
		if v {
			s = append(s, t.symbols[uint32(i)].Name())
		}
	}
	return s
}

// ****************************************************************************

// Bitfield32 represents the sunspec type bitfield32.
type Bitfield32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint32) error
	// Get returns the point´s underlying value.
	Get() uint32
	// Flip sets the bit at position pos, starting at 0, to the value of v.
	Flip(pos int, v bool) error
	// Field returns the individual bit values as bool array.
	Field() [32]bool
	// States returns all active enumerated states, correlating the bit value to its symbol.
	States() []string
}

type tBitfield32 struct {
	point
	data    uint32
	symbols Symbols
}

var _ Bitfield32 = (*tBitfield32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tBitfield32) Valid() bool { return t.Get() != 0xFFFFFFFF }

// String formats the point´s value as string.
func (t *tBitfield32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tBitfield32) Quantity() uint16 { return 2 }

// encode puts the point´s value into a buffer.
func (t *tBitfield32) encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tBitfield32) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint32(buf))
}

// Set sets the point´s underlying value.
func (t *tBitfield32) Set(v uint32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tBitfield32) Get() uint32 { return t.data }

// Flip sets the bit at position pos, starting at 0, to the value of v.
func (t *tBitfield32) Flip(pos int, v bool) error {
	switch {
	case pos < 0 || pos > 31:
		return errors.New("sunspec: out of bounds bit position ")
	case v:
		return t.Set(t.Get() | (1 << pos))
	}
	return t.Set(t.Get() &^ (1 << pos))
}

// Field returns the individual bit values as bool array.
func (t *tBitfield32) Field() (f [32]bool) {
	for v, b := t.Get(), 0; b < len(f); b++ {
		f[b] = v&(1<<b) != 0
	}
	return f
}

// States returns all active enumerated states, correlating the bit value to its symbol.
func (t *tBitfield32) States() (s []string) {
	if !t.Valid() {
		return nil
	}
	for i, v := range t.Field() {
		if v {
			s = append(s, t.symbols[uint32(i)].Name())
		}
	}
	return s
}

// ****************************************************************************

// Bitfield64 represents the sunspec type bitfield64.
type Bitfield64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint64) error
	// Get returns the point´s underlying value.
	Get() uint64
	// Flip sets the bit at position pos, starting at 0, to the value of v.
	Flip(pos int, v bool) error
	// Field returns the individual bit values as bool array.
	Field() [64]bool
	// States returns all active enumerated states, correlating the bit value to its symbol.
	States() []string
}

type tBitfield64 struct {
	point
	data    uint64
	symbols Symbols
}

var _ Bitfield64 = (*tBitfield64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tBitfield64) Valid() bool { return t.Get() != 0xFFFFFFFFFFFFFFFF }

// String formats the point´s value as string.
func (t *tBitfield64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tBitfield64) Quantity() uint16 { return 4 }

// encode puts the point´s value into a buffer.
func (t *tBitfield64) encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tBitfield64) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint64(buf))
}

// Set sets the point´s underlying value.
func (t *tBitfield64) Set(v uint64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tBitfield64) Get() uint64 { return t.data }

// Flip sets the bit at position pos, starting at 0, to the value of v.
func (t *tBitfield64) Flip(pos int, v bool) error {
	switch {
	case pos < 0 || pos > 63:
		return errors.New("sunspec: out of bounds bit position ")
	case v:
		return t.Set(t.Get() | (1 << pos))
	}
	return t.Set(t.Get() &^ (1 << pos))
}

// Field returns the individual bit values as bool array.
func (t *tBitfield64) Field() (f [64]bool) {
	for v, b := t.Get(), 0; b < len(f); b++ {
		f[b] = v&(1<<b) != 0
	}
	return f
}

// States returns all active enumerated states, correlating the bit value to its symbol.
func (t *tBitfield64) States() (s []string) {
	if !t.Valid() {
		return nil
	}
	for i, v := range t.Field() {
		if v {
			s = append(s, t.symbols[uint32(i)].Name())
		}
	}
	return s
}

// ****************************************************************************

// Enum16 represents the sunspec type enum16.
type Enum16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint16) error
	// Get returns the point´s underlying value.
	Get() uint16
	// State returns the currently active enumerated state.
	State() string
}

type tEnum16 struct {
	point
	data    uint16
	symbols Symbols
}

var _ Enum16 = (*tEnum16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tEnum16) Valid() bool { return t.Get() != 0xFFFF }

// String formats the point´s value as string.
func (t *tEnum16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tEnum16) Quantity() uint16 { return 1 }

// encode puts the point´s value into a buffer.
func (t *tEnum16) encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tEnum16) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *tEnum16) Set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tEnum16) Get() uint16 { return t.data }

// State returns the currently active enumerated state.
func (t *tEnum16) State() string { return t.symbols[uint32(t.Get())].Name() }

// ****************************************************************************

// Enum32 represents the sunspec type enum32.
type Enum32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint32) error
	// Get returns the point´s underlying value.
	Get() uint32
	// State returns the currently active enumerated state.
	State() string
}

type tEnum32 struct {
	point
	data    uint32
	symbols Symbols
}

var _ Enum32 = (*tEnum32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tEnum32) Valid() bool { return t.Get() != 0xFFFFFFFF }

// String formats the point´s value as string.
func (t *tEnum32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tEnum32) Quantity() uint16 { return 2 }

// encode puts the point´s value into a buffer.
func (t *tEnum32) encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tEnum32) decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint32(buf))
}

// Set sets the point´s underlying value.
func (t *tEnum32) Set(v uint32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tEnum32) Get() uint32 { return t.data }

// State returns the currently active enumerated state.
func (t *tEnum32) State() string { return t.symbols[t.Get()].Name() }

// ****************************************************************************

// String represents the sunspec type string.
type String interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v string) error
	// Get returns the point´s underlying value.
	Get() string
}

type tString struct {
	point
	data []byte
}

var _ String = (*tString)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tString) Valid() bool { return t.Get() != "" }

// String formats the point´s value as string.
func (t *tString) String() string { return t.Get() }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tString) Quantity() uint16 { return uint16(cap(t.data) / 2) }

// encode puts the point´s value into a buffer.
func (t *tString) encode(buf []byte) error {
	copy(buf, []byte(t.Get()))
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tString) decode(buf []byte) error {
	return t.Set(string(buf[:2*t.Quantity()]))
}

// Set sets the point´s underlying value.
func (t *tString) Set(v string) error {
	copy(t.data[:cap(t.data)], v)
	return nil
}

// Get returns the point´s underlying value.
func (t *tString) Get() string { return string(t.data) }

// ****************************************************************************

// Float32 represents the sunspec type float32.
type Float32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v float32) error
	// Get returns the point´s underlying value.
	Get() float32
}

type tFloat32 struct {
	point
	data float32
}

var _ Float32 = (*tFloat32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tFloat32) Valid() bool { return t.Get() != 0x7FC00000 }

// String formats the point´s value as string.
func (t *tFloat32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tFloat32) Quantity() uint16 { return 2 }

// encode puts the point´s value into a buffer.
func (t *tFloat32) encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, math.Float32bits(t.Get()))
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tFloat32) decode(buf []byte) error {
	return t.Set(math.Float32frombits(binary.BigEndian.Uint32(buf)))
}

// Set sets the point´s underlying value.
func (t *tFloat32) Set(v float32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tFloat32) Get() float32 { return t.data }

// ****************************************************************************

// Float64 represents the sunspec type float64.
type Float64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v float64) error
	// Get returns the point´s underlying value.
	Get() float64
}

type tFloat64 struct {
	point
	data float64
}

var _ Float64 = (*tFloat64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tFloat64) Valid() bool { return t.Get() != 0x7FC00000 }

// String formats the point´s value as string.
func (t *tFloat64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tFloat64) Quantity() uint16 { return 4 }

// encode puts the point´s value into a buffer.
func (t *tFloat64) encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, math.Float64bits(t.Get()))
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tFloat64) decode(buf []byte) error {
	return t.Set(math.Float64frombits(binary.BigEndian.Uint64(buf)))
}

// Set sets the point´s underlying value.
func (t *tFloat64) Set(v float64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *tFloat64) Get() float64 { return t.data }

// ****************************************************************************

// Ipaddr represents the sunspec type ipaddr.
type Ipaddr interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v net.IP) error
	// Get returns the point´s underlying value.
	Get() net.IP
	// Raw returns the point´s raw data.
	Raw() [4]byte
}

type tIpaddr struct {
	point
	data [4]byte
}

var _ Ipaddr = (*tIpaddr)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tIpaddr) Valid() bool { return t.data != [4]byte{} }

// String formats the point´s value as string.
func (t *tIpaddr) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tIpaddr) Quantity() uint16 { return uint16(len(t.data) / 2) }

// encode puts the point´s value into a buffer.
func (t *tIpaddr) encode(buf []byte) error {
	copy(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tIpaddr) decode(buf []byte) error {
	return t.Set(buf)
}

// Set sets the point´s underlying value.
func (t *tIpaddr) Set(v net.IP) error {
	copy(t.data[:len(t.data)], v)
	return nil
}

// Get returns the point´s underlying value.
func (t *tIpaddr) Get() net.IP { return append(net.IP(nil), t.data[:]...) }

// Raw returns the point´s raw data.
func (t *tIpaddr) Raw() (r [4]byte) {
	copy(r[:], t.Get())
	return r
}

// ****************************************************************************

// Ipaddr represents the sunspec type ipaddr.
type Ipv6addr interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v net.IP) error
	// Get returns the point´s underlying value.
	Get() net.IP
	// Raw returns the point´s raw data.
	Raw() [16]byte
}

type tIpv6addr struct {
	point
	data [16]byte
}

var _ Ipv6addr = (*tIpv6addr)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tIpv6addr) Valid() bool { return t.data != [16]byte{} }

// String formats the point´s value as string.
func (t *tIpv6addr) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tIpv6addr) Quantity() uint16 { return uint16(len(t.data) / 2) }

// encode puts the point´s value into a buffer.
func (t *tIpv6addr) encode(buf []byte) error {
	copy(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tIpv6addr) decode(buf []byte) error {
	return t.Set(buf)
}

// Set sets the point´s underlying value.
func (t *tIpv6addr) Set(v net.IP) error {
	copy(t.data[:len(t.data)], v)
	return nil
}

// Get returns the point´s underlying value.
func (t *tIpv6addr) Get() net.IP { return append(net.IP(nil), t.data[:]...) }

// Raw returns the point´s raw data.
func (t *tIpv6addr) Raw() (r [16]byte) {
	copy(r[:], t.Get())
	return r
}

// ****************************************************************************

// Eui48 represents the sunspec type eui48.
type Eui48 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v net.HardwareAddr) error
	// Get returns the point´s underlying value.
	Get() net.HardwareAddr
	// Raw returns the point´s raw data.
	Raw() [8]byte
}

type tEui48 struct {
	point
	data [8]byte
}

var _ Eui48 = (*tEui48)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *tEui48) Valid() bool { return true } //?

// String formats the point´s value as string.
func (t *tEui48) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *tEui48) Quantity() uint16 { return uint16(len(t.data) / 2) }

// encode puts the point´s value into a buffer.
func (t *tEui48) encode(buf []byte) error {
	copy(buf, t.Get())
	return nil
}

// decode sets the point´s value from a buffer.
func (t *tEui48) decode(buf []byte) error {
	return t.Set(buf)
}

// Set sets the point´s underlying value.
func (t *tEui48) Set(v net.HardwareAddr) error {
	copy(t.data[:len(t.data)], v)
	return nil
}

// Get returns the point´s underlying value.
func (t *tEui48) Get() net.HardwareAddr { return append(net.HardwareAddr(nil), t.data[:]...) }

// Raw returns the point´s raw data.
func (t *tEui48) Raw() (r [8]byte) {
	copy(r[:], t.Get())
	return r
}
