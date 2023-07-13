package goxdr

import (
	"fmt"
	"math"
	"errors"
)

type PrimitiveReadState struct {
	primitiveSize int
	bytes [8]byte
	fillCount int
}

func NewPrimitiveReadState(primitiveSize int) (state *PrimitiveReadState, err error) {
	switch primitiveSize {
		case 4, 8:
			state = &PrimitiveReadState {
				primitiveSize: primitiveSize,
			}
		default:
			err = errors.New(fmt.Sprintf("Expected primitive size to be 4 or 8, not %d", primitiveSize))
	}
	return
}

func(state *PrimitiveReadState) Reset(primitiveSize int) error {
	switch primitiveSize {
		case 4, 8:
			state.primitiveSize = primitiveSize
			state.fillCount = 0
			return nil
		default:
			return errors.New(fmt.Sprintf("Expected primitive size to be 4 or 8, not %d", primitiveSize))
	}
}

func(state *PrimitiveReadState) Update(bytes []byte) (readCount int, isFull bool) {
	if state.fillCount > state.primitiveSize {
		panic(fmt.Sprintf("fillCount (%d) > primitiveSize (%d)", state.fillCount, state.primitiveSize))
	}
	need := len(bytes)
	can := state.primitiveSize - state.fillCount
	if can < need {
		readCount = can
	} else {
		readCount = need
	}
	for index := 0; index < readCount; index++ {
		state.bytes[state.fillCount + index] = bytes[index]
	}
	state.fillCount += readCount
	if state.fillCount == state.primitiveSize {
		isFull = true
	}
	return
}

func(state *PrimitiveReadState) EndPacket() error {
	if state.fillCount > state.primitiveSize {
		panic(fmt.Sprintf("fillCount (%d) > primitiveSize (%d)", state.fillCount, state.primitiveSize))
	}
	if state.fillCount == state.primitiveSize {
		return nil
	}
	return errors.New(fmt.Sprintf(
		"Missing %d bytes for primitive of size %d",
		state.primitiveSize - state.fillCount,
		state.primitiveSize,
	))
}

func(state *PrimitiveReadState) AsInt() int32 {
	var digits uint32 = (uint32(state.bytes[0]) << 24) |
		(uint32(state.bytes[1]) << 16) |
		(uint32(state.bytes[2]) << 8) |
		uint32(state.bytes[3])
	return int32(digits)
}

func(state *PrimitiveReadState) AsUint() uint32 {
	return (uint32(state.bytes[0]) << 24) |
		(uint32(state.bytes[1]) << 16) |
		(uint32(state.bytes[2]) << 8) |
		uint32(state.bytes[3])
}

func(state *PrimitiveReadState) AsHyperInt() int64 {
	var digits uint64 = (uint64(state.bytes[0]) << 56) |
		(uint64(state.bytes[1]) << 48) |
		(uint64(state.bytes[2]) << 40) |
		(uint64(state.bytes[3]) << 32) |
		(uint64(state.bytes[4]) << 24) |
		(uint64(state.bytes[5]) << 16) |
		(uint64(state.bytes[6]) << 8) |
		uint64(state.bytes[7])
	return int64(digits)
}

func(state *PrimitiveReadState) AsHyperUint() uint64 {
	return (uint64(state.bytes[0]) << 56) |
		(uint64(state.bytes[1]) << 48) |
		(uint64(state.bytes[2]) << 40) |
		(uint64(state.bytes[3]) << 32) |
		(uint64(state.bytes[4]) << 24) |
		(uint64(state.bytes[5]) << 16) |
		(uint64(state.bytes[6]) << 8) |
		uint64(state.bytes[7])
}

func(state *PrimitiveReadState) AsFloat() float32 {
	return math.Float32frombits(state.AsUint())
}

func(state *PrimitiveReadState) AsDouble() float64 {
	return math.Float64frombits(state.AsHyperUint())
}

var _ ReadState = &PrimitiveReadState{}
