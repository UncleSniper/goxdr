package goxdr

import (
	"io"
	"math"
)

func WriteInt(value int32, buffer []byte, writer io.Writer) (err error) {
	digits := uint32(value)
	buffer[0] = byte(digits >> 24)
	buffer[1] = byte(digits >> 16)
	buffer[2] = byte(digits >> 8)
	buffer[3] = byte(digits)
	_, err = writer.Write(buffer[0:4])
	return
}

func WriteUint(value uint32, buffer []byte, writer io.Writer) (err error) {
	buffer[0] = byte(value >> 24)
	buffer[1] = byte(value >> 16)
	buffer[2] = byte(value >> 8)
	buffer[3] = byte(value)
	_, err = writer.Write(buffer[0:4])
	return
}

func WriteHyperInt(value int64, buffer []byte, writer io.Writer) (err error) {
	digits := uint64(value)
	buffer[0] = byte(digits >> 56)
	buffer[1] = byte(digits >> 48)
	buffer[2] = byte(digits >> 40)
	buffer[3] = byte(digits >> 32)
	buffer[4] = byte(digits >> 24)
	buffer[5] = byte(digits >> 16)
	buffer[6] = byte(digits >> 8)
	buffer[7] = byte(digits)
	_, err = writer.Write(buffer[0:8])
	return
}

func WriteHyperUint(value uint64, buffer []byte, writer io.Writer) (err error) {
	buffer[0] = byte(value >> 56)
	buffer[1] = byte(value >> 48)
	buffer[2] = byte(value >> 40)
	buffer[3] = byte(value >> 32)
	buffer[4] = byte(value >> 24)
	buffer[5] = byte(value >> 16)
	buffer[6] = byte(value >> 8)
	buffer[7] = byte(value)
	_, err = writer.Write(buffer[0:8])
	return
}

func WriteFloat(value float32, buffer []byte, writer io.Writer) error {
	return WriteUint(math.Float32bits(value), buffer, writer)
}

func WriteDouble(value float64, buffer []byte, writer io.Writer) error {
	return WriteHyperUint(math.Float64bits(value), buffer, writer)
}
