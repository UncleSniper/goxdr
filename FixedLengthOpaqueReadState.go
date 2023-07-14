package goxdr

import (
	"fmt"
	"math"
	"errors"
)

type FixedLengthOpaqueReadState struct {
	ExpectedLength uint32
	Handler ReadState
	HandlerName string
	currentLength uint32
	firstError error
}

func(state *FixedLengthOpaqueReadState) Reset() {
	state.currentLength = 0
	state.firstError = nil
}

func(state *FixedLengthOpaqueReadState) Update(bytes []byte) (readCount int, isFull bool) {
	if state.firstError != nil {
		isFull = true
		return
	}
	length := len(bytes)
	if int64(length) > int64(math.MaxUint32) {
		state.firstError = errors.New(fmt.Sprintf(
			"Update with chunk of size %d, which exceeds the range of uint32",
			length,
		))
		isFull = true
		return
	}
	length32 := uint32(length)
	nextLength := state.currentLength + length32
	if nextLength < length32 {
		state.firstError = errors.New(fmt.Sprintf(
			"Update with chunk of size %d would cause total stream size (%d before chunk) " +
					"to exceed the range of uint32",
			length32,
			state.currentLength,
		))
		isFull = true
		return
	}
	if nextLength > state.ExpectedLength {
		length32 = state.ExpectedLength - state.currentLength
	}
	if length32 > 0 {
		readCount, isFull = state.Handler.Update(bytes[0:length32])
		if uint32(readCount) > length32 {
			state.firstError = errors.New(fmt.Sprintf(
				"Opaque data handler read %d bytes, but was supposed to only read %d",
				readCount,
				length32,
			))
			isFull = true
			return
		}
		state.currentLength += uint32(readCount)
	} else if state.currentLength == state.ExpectedLength {
		isFull = true
	}
	return
}

func(state *FixedLengthOpaqueReadState) EndPacket() (err error) {
	if state.firstError != nil {
		err = state.firstError
	} else {
		err = state.Handler.EndPacket()
		if err != nil {
			err = &OpaqueHandlerError {
				PropagatedError: err,
				HandlerName: state.HandlerName,
			}
		}
	}
	return
}

var _ ReadState = &FixedLengthOpaqueReadState{}
