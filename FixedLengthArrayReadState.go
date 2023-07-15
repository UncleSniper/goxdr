package goxdr

import (
	"fmt"
	"errors"
)

type FixedLengthArrayReadState[T any] struct {
	ExpectedLength uint32
	HandlerFactory TypedReadStateFactory[T]
	HandlerName string
	currentIndex uint32
	currentHandler ReadState
	firstError error
}

func(state *FixedLengthArrayReadState[T]) Reset() {
	state.currentIndex = 0
	state.currentHandler = nil
	state.firstError = nil
}

func(state *FixedLengthArrayReadState[T]) nextHandler() bool {
	state.currentHandler, state.firstError = state.HandlerFactory(state.currentIndex, state.ExpectedLength)
	if state.firstError != nil {
		return true
	}
	if state.currentHandler != nil {
		return false
	}
	state.firstError = errors.New(fmt.Sprintf(
		"Read state factory returned nil for index %d of %d",
		state.currentIndex,
		state.ExpectedLength,
	))
	return true
}

func(state *FixedLengthArrayReadState[T]) Update(bytes []byte) (readCount int, isFull bool) {
	if state.firstError != nil {
		isFull = true
		return
	}
	if state.currentIndex >= state.ExpectedLength {
		isFull = true
		return
	}
	if state.currentHandler == nil {
		isFull = state.nextHandler()
		if isFull {
			return
		}
	}
	var handled int
	for {
		handled, isFull = state.currentHandler.Update(bytes[readCount:])
		readCount += handled
		if !isFull {
			return
		}
		state.firstError = state.currentHandler.EndPacket()
		if state.firstError != nil {
			isFull = true
			return
		}
		state.currentIndex++
		if state.currentIndex < state.ExpectedLength {
			isFull = state.nextHandler()
			if isFull {
				return
			}
		} else {
			state.currentHandler = nil
			isFull = true
			return
		}
	}
}

func(state *FixedLengthArrayReadState[T]) EndPacket() (err error) {
	if state.firstError == nil && state.currentIndex < state.ExpectedLength {
		for {
			state.firstError = state.currentHandler.EndPacket()
			if state.firstError != nil {
				break
			}
			state.currentIndex++
			if state.currentIndex >= state.ExpectedLength {
				break
			}
			if state.nextHandler() {
				break
			}
		}
	}
	return state.firstError
}

var _ ReadState = &FixedLengthArrayReadState[int]{}
