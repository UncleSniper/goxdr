package goxdr

import (
	"fmt"
	"errors"
)

type VariableLengthArrayReadState[T any] struct {
	PrimitiveState *PrimitiveReadState
	FixedLengthState *FixedLengthArrayReadState[T]
	MaxLength uint32
	inBody bool
	firstError error
}

func(state *VariableLengthArrayReadState[T]) Reset() {
	state.PrimitiveState.Reset(4)
	state.inBody = false
	state.firstError = nil
}

func(state *VariableLengthArrayReadState[T]) Update(bytes []byte) (readCount int, isFull bool) {
	if state.firstError != nil {
		isFull = true
		return
	}
	length := len(bytes)
	if !state.inBody {
		readCount, isFull = state.PrimitiveState.Update(bytes)
		if readCount > length {
			state.firstError = errors.New(fmt.Sprintf(
				"Primitive read state read %d bytes, but was supposed to only read %d",
				readCount,
				length,
			))
			isFull = true
			return
		}
		if !isFull {
			return
		}
		state.firstError = state.PrimitiveState.EndPacket()
		if state.firstError != nil {
			return
		}
		state.FixedLengthState.Reset()
		state.FixedLengthState.ExpectedLength = state.PrimitiveState.AsUint()
		if state.FixedLengthState.ExpectedLength > state.MaxLength {
			state.firstError = errors.New(fmt.Sprintf(
				"Variable-length array has maximum length %d, but encountered length %d",
				state.MaxLength,
				state.FixedLengthState.ExpectedLength,
			))
			return
		}
		state.inBody = true
	}
	var bodyReadCount int
	bodyReadCount, isFull = state.FixedLengthState.Update(bytes[readCount:])
	if bodyReadCount > length - readCount {
		state.firstError = errors.New(fmt.Sprintf(
			"Fixed length array read state read %d bytes, but was supposed to only read %d",
			bodyReadCount,
			length - readCount,
		))
		isFull = true
		return
	}
	readCount += bodyReadCount
	return
}

func(state *VariableLengthArrayReadState[T]) EndPacket() error {
	if state.firstError == nil {
		if state.inBody {
			state.firstError = state.FixedLengthState.EndPacket()
		} else {
			state.firstError = state.PrimitiveState.EndPacket()
			if state.firstError == nil {
				state.FixedLengthState.Reset()
				state.FixedLengthState.ExpectedLength = state.PrimitiveState.AsUint()
				state.inBody = true
				state.firstError = state.FixedLengthState.EndPacket()
			}
		}
	}
	return state.firstError
}

var _ ReadState = &VariableLengthArrayReadState[int]{}
