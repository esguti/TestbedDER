package modbus

import (
	"encoding/binary"
)

func boundCheck(address, quantity, limit uint16) Exception {
	switch {
	case quantity < 1 || quantity > limit:
		return IllegalDataValue
	case address > address+quantity-1:
		return IllegalDataAddress
	}
	return 0
}

func byteCount(bitCount uint16) int {
	return int((bitCount + 7) / 8)
}

func bytesToBools(quantity uint16, bytes []byte) []bool {
	buf := make([]bool, quantity)
	for i, x := range bytes {
		for j := 0; j < 8; j++ {
			k := 8*i + j
			if len(buf) == k {
				return buf
			}
			buf[k] = (x<<uint(j))&0x80 == 0x80
		}
	}
	return buf
}

func put(length int, args ...interface{}) []byte {
	new := make([]byte, length)
	buf := new
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			buf = putBool(buf, v)
		case []bool:
			buf = putBoolS(buf, v)
		case byte:
			buf = putByte(buf, v)
		case []byte:
			buf = putByteS(buf, v)
		case uint16:
			buf = putUint16(buf, v)
		case []uint16:
			buf = putUint16S(buf, v)
		}
	}

	return new
}

func putBool(buf []byte, arg bool) []byte {
	if arg {
		return putUint16(buf, 0xFF00)
	}
	return putUint16(buf, 0x0000)
}

func putBoolS(buf []byte, args []bool) []byte {
	var x bool
	var c int
	for c, x = range args {
		if x {
			buf[c/8] |= 0x80 >> uint(c%8)
		}
	}
	return buf[c/8+1:]
}

func putByte(buf []byte, arg byte) []byte {
	buf[0] = arg
	return buf[1:]
}

func putByteS(buf []byte, args []byte) []byte {
	return buf[copy(buf, args):]
}

func putUint16(buf []byte, arg uint16) []byte {
	binary.BigEndian.PutUint16(buf, arg)
	return buf[2:]
}

func putUint16S(buf []byte, args []uint16) []byte {
	for _, arg := range args {
		buf = putUint16(buf, arg)
	}
	return buf
}
