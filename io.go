package goxdr

import (
	"io"
	"fmt"
	"math"
	"errors"
)

type countingWriter struct {
	writer io.Writer
	count uint32
}

func(counter *countingWriter) Write(bytes []byte) (writeCount int, err error) {
	writeCount, err = counter.writer.Write(bytes)
	if err == nil {
		if int64(writeCount) > int64(math.MaxUint32) {
			err = errors.New(fmt.Sprintf("Wrote chunk of size %d, which exceeds the range of uint32", writeCount))
			return
		}
		nextCount := counter.count + uint32(writeCount)
		if nextCount < counter.count {
			err = errors.New(fmt.Sprintf(
				"Chunk (of size %d) would cause total write size (%d before chunk) to exceed the range of uint32",
				writeCount,
				counter.count,
			))
			return
		}
		counter.count = nextCount
	}
	return
}

var _ io.Writer = &countingWriter{}
