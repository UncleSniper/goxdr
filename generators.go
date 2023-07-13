package goxdr

import (
	"io"
)

type ByteGenerator func([]byte, io.Writer) error

type BytePadder func(io.Writer, uint32) error

var zeroBytesSlice []byte = make([]byte, zeroSliceSize)

func ZeroBytePadder(writer io.Writer, count uint32) (err error) {
	whole := count / uint32(zeroSliceSize)
	rest := int(count % uint32(zeroSliceSize))
	for u := uint32(0); u < whole; u++ {
		_, err = writer.Write(zeroBytesSlice)
		if err != nil {
			return
		}
	}
	if rest > 0 {
		_, err = writer.Write(zeroBytesSlice[0:rest])
	}
	return
}

type PacketSink[T any] func(TypedPacket[T]) error

type PacketGenerator[T any] func(PacketSink[T]) error

type ElementPadder[T any] func([]byte, io.Writer) error

func TypedPacketSliceGenerator[T any](elements []TypedPacket[T]) PacketGenerator[T] {
	return func(sink PacketSink[T]) (err error) {
		for _, element := range elements {
			err = sink(element)
			if err != nil {
				return
			}
		}
		return
	}
}
