package goxdr

import (
	"errors"
)

type ReadState interface {
	Update([]byte) (int, bool)
	EndPacket() error
}

type RequestReadState interface {
	ReadState
	ResponsePacket() Packet
}

type ReadStateFactory func(uint32, uint32) (ReadState, error)

type TypedReadState[T any] interface {
	ReadState
}

type TypedReadStateFactory[T any] func(uint32, uint32) (TypedReadState[T], error)

func TypedReadStateFactoryOf[T any](factory ReadStateFactory) TypedReadStateFactory[T] {
	return func(index uint32, size uint32) (adapted TypedReadState[T], err error) {
		if factory == nil {
			err = errors.New("Read state factory is nil")
		} else {
			adapted, err = factory(index, size)
		}
		return
	}
}
