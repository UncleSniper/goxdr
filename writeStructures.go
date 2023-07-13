package goxdr

import (
	"io"
	"fmt"
	"math"
	"errors"
)

func writeRemainder(buffer []byte, writer io.Writer, remainder int) (err error) {
	remainder = 4 - remainder
	for i := 0; i < remainder; i++ {
		buffer[i] = 0
	}
	_, err = writer.Write(buffer[0:remainder])
	return
}

func WriteFixedLengthOpaquePacket(packet Packet, buffer []byte, writer io.Writer) (err error) {
	remainder := int(packet.ByteSize() % uint32(4))
	err = packet.WriteTo(buffer, writer)
	if err != nil || remainder == 0 {
		return
	}
	err = writeRemainder(buffer, writer, remainder)
	return
}

func WriteFixedLengthOpaqueReader(
	reader io.Reader,
	expectedSize uint32,
	buffer []byte,
	writer io.Writer,
	padding BytePadder,
) (err error) {
	var transferBuffer []byte
	if len(buffer) >= minBulkTransferBufferSize {
		transferBuffer = buffer
	} else {
		transferBuffer = make([]byte, minBulkTransferBufferSize)
	}
	var actualSize uint32
	var readCount int
	for {
		readCount, err = reader.Read(transferBuffer)
		if int64(readCount) > int64(math.MaxUint32) {
			err = errors.New(fmt.Sprintf("Read chunk of size %d, which exceeds the range of uint32", readCount))
			return
		}
		nextSize := actualSize + uint32(readCount)
		if nextSize < actualSize {
			err = errors.New(fmt.Sprintf(
				"Chunk (of size %d) would cause total read size (%d before chunk) to exceed the range of uint32",
				readCount,
				actualSize,
			))
			return
		}
		var writeErr error
		if readCount > 0 {
			_, writeErr = writer.Write(transferBuffer[0:readCount])
		}
		eof := err == io.EOF
		if err == nil || eof {
			err = writeErr
		}
		actualSize = nextSize
		if eof || err != nil {
			break
		}
	}
	if err == nil && actualSize != expectedSize {
		if padding != nil && actualSize < expectedSize {
			err = padding(writer, expectedSize - actualSize)
		} else {
			err = errors.New(fmt.Sprintf(
				"Expected stream length (%d) does not match actual stream length (%d)",
				expectedSize,
				actualSize,
			))
		}
	}
	remainder := int(expectedSize % uint32(4))
	if err == nil && remainder > 0 {
		err = writeRemainder(buffer, writer, remainder)
	}
	return
}

func WriteFixedLengthOpaqueGenerator(
	generator ByteGenerator,
	expectedSize uint32,
	buffer []byte,
	writer io.Writer,
	padding BytePadder,
) (err error) {
	var counter countingWriter
	counter.writer = writer
	err = generator(buffer, &counter)
	if err == nil && counter.count != expectedSize {
		if padding != nil && counter.count < expectedSize {
			err = padding(writer, expectedSize - counter.count)
		} else {
			err = errors.New(fmt.Sprintf(
				"Expected stream length (%d) does not match actual stream length (%d)",
				expectedSize,
				counter.count,
			))
		}
	}
	remainder := int(expectedSize % uint32(4))
	if err == nil && remainder > 0 {
		err = writeRemainder(buffer, writer, remainder)
	}
	return
}

func WriteVariableLengthOpaquePacket(packet Packet, maxSize uint32, buffer []byte, writer io.Writer) (err error) {
	actualSize := packet.ByteSize()
	if actualSize > maxSize {
		err = errors.New(fmt.Sprintf("Packet length (%d) exceeds maximum length (%d)", actualSize, maxSize))
		return
	}
	err = WriteUint(actualSize, buffer, writer)
	if err != nil {
		return
	}
	remainder := int(actualSize % uint32(4))
	err = packet.WriteTo(buffer, writer)
	if err == nil && remainder > 0 {
		err = writeRemainder(buffer, writer, remainder)
	}
	return
}

func WriteVariableLengthOpaqueReader(
	reader io.Reader,
	expectedSize uint32,
	maxSize uint32,
	buffer []byte,
	writer io.Writer,
	padding BytePadder,
) (err error) {
	if expectedSize > maxSize {
		err = errors.New(fmt.Sprintf("Packet length (%d) exceeds maximum length (%d)", expectedSize, maxSize))
		return
	}
	err = WriteUint(expectedSize, buffer, writer)
	if err == nil {
		err = WriteFixedLengthOpaqueReader(reader, expectedSize, buffer, writer, padding)
	}
	return
}

func WriteVariableLengthOpaqueGenerator(
	generator ByteGenerator,
	expectedSize uint32,
	maxSize uint32,
	buffer []byte,
	writer io.Writer,
	padding BytePadder,
) (err error) {
	if expectedSize > maxSize {
		err = errors.New(fmt.Sprintf("Packet length (%d) exceeds maximum length (%d)", expectedSize, maxSize))
		return
	}
	err = WriteUint(expectedSize, buffer, writer)
	if err == nil {
		WriteFixedLengthOpaqueGenerator(generator, expectedSize, buffer, writer, padding)
	}
	return
}

func WriteFixedLengthArrayGenerator[T any](
	generator PacketGenerator[T],
	expectedSize uint32,
	buffer []byte,
	writer io.Writer,
	padding ElementPadder[T],
) (err error) {
	var actualSize uint32
	err = generator(func(packet TypedPacket[T]) error {
		actualSize++
		if actualSize == 0 {
			return errors.New("Generated element caused actual element count to exceed the range of uint32")
		} else {
			return packet.WriteTo(buffer, writer)
		}
	})
	if err == nil && actualSize != expectedSize {
		if padding != nil && actualSize < expectedSize {
			for u := actualSize; u < expectedSize; u++ {
				err = padding(buffer, writer)
				if err != nil {
					return
				}
			}
		} else {
			err = errors.New(fmt.Sprintf(
				"Expected array length (%d) does not match actual array length (%d)",
				expectedSize,
				actualSize,
			))
		}
	}
	return
}

func WriteVariableLengthArrayGenerator[T any](
	generator PacketGenerator[T],
	expectedSize uint32,
	maxSize uint32,
	buffer []byte,
	writer io.Writer,
	padding ElementPadder[T],
) (err error) {
	if expectedSize > maxSize {
		err = errors.New(fmt.Sprintf(
			"Expected element count (%d) exceeds maximum count (%d)",
			expectedSize,
			maxSize,
		))
		return
	}
	err = WriteUint(expectedSize, buffer, writer)
	if err == nil {
		err = WriteFixedLengthArrayGenerator(generator, expectedSize, buffer, writer, padding)
	}
	return
}
