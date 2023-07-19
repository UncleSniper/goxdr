package goxdr

import (
	"fmt"
	"errors"
)

type TaggedUnionReadState[T any] struct {
	PrimitiveState *PrimitiveReadState
	HandlerFactory TypedReadStateFactory[T]
	HandlerName string
	currentHandler ReadState
	firstError error
}

func(state *TaggedUnionReadState[T]) Reset() {
	state.PrimitiveState.Reset(4)
	state.currentHandler = nil
	state.firstError = nil
}

func(state *TaggedUnionReadState[T]) enterArm() bool {
	discriminant := state.PrimitiveState.AsUint()
	state.currentHandler, state.firstError = state.HandlerFactory(discriminant, 0)
	if state.firstError != nil {
		return true
	}
	if state.currentHandler == nil {
		state.firstError = &UnionDiscriminantError {
			Discriminant: discriminant,
			HandlerName: state.HandlerName,
		}
		return true
	}
	return false
}

func(state *TaggedUnionReadState[T]) Update(bytes []byte) (readCount int, isFull bool) {
	if state.firstError != nil {
		isFull = true
		return
	}
	length := len(bytes)
	if state.currentHandler == nil {
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
		if state.firstError != nil || state.enterArm() {
			return
		}
	}
	var armReadCount int
	armReadCount, isFull = state.currentHandler.Update(bytes[readCount:])
	if armReadCount > length - readCount {
		state.firstError = errors.New(fmt.Sprintf(
			"Tagged union arm read state read %d bytes, but was supposed to only read %d",
			armReadCount,
			length - readCount,
		))
		isFull = true
		return
	}
	readCount += armReadCount
	return
}

func(state *TaggedUnionReadState[T]) EndPacket() error {
	if state.firstError == nil {
		if state.currentHandler != nil {
			state.firstError = state.currentHandler.EndPacket()
		} else {
			state.firstError = state.PrimitiveState.EndPacket()
			if state.firstError == nil && !state.enterArm() {
				state.firstError = state.currentHandler.EndPacket()
			}
		}
	}
	return state.firstError
}

var _ ReadState = &TaggedUnionReadState[int]{}
