package goxdr

import (
	"io"
	"fmt"
	"math"
	"errors"
)

type Packet interface {
	ByteSize() uint32
	WriteTo([]byte, io.Writer) error
}

type TypedPacket[T any] interface {
	Packet
}

type ByteSlicePacket struct {
	Bytes []byte
}

func(packet ByteSlicePacket) ByteSize() uint32 {
	size := len(packet.Bytes)
	if int64(size) > int64(math.MaxUint32) {
		panic(fmt.Sprintf("Size of slice (%d elements) exceeds range of uint32", size))
	}
	return uint32(size)
}

func(packet ByteSlicePacket) WriteTo(buffer []byte, writer io.Writer) (err error) {
	if len(packet.Bytes) > 0 {
		_, err = writer.Write(packet.Bytes)
	}
	return
}

type PaddingPacket struct {
	ShortPacket Packet
	RequiredLength uint32
	Padding BytePadder
}

func(packet *PaddingPacket) ByteSize() uint32 {
	return packet.RequiredLength
}

func(packet *PaddingPacket) WriteTo(buffer []byte, writer io.Writer) (err error) {
	shortLength := packet.ShortPacket.ByteSize()
	if shortLength > packet.RequiredLength {
		err = errors.New(fmt.Sprintf(
			"Cannot pad Packet of size %d to size %d: Padding would be negative length",
			shortLength,
			packet.RequiredLength,
		))
	}
	err = packet.ShortPacket.WriteTo(buffer, writer)
	if err != nil {
		return
	}
	remainder := packet.RequiredLength - shortLength
	if remainder > 0 {
		padding := packet.Padding
		if padding == nil {
			padding = ZeroBytePadder
		}
		err = padding(writer, remainder)
	}
	return
}

var _ Packet = ByteSlicePacket{}
var _ Packet = &PaddingPacket{}
