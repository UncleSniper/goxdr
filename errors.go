package goxdr

import (
	"strings"
	"strconv"
)

type OpaqueHandlerError struct {
	PropagatedError error
	HandlerName string
}

func(err *OpaqueHandlerError) Error() string {
	var builder strings.Builder
	if len(err.HandlerName) > 0 {
		builder.WriteString(err.HandlerName)
		builder.WriteString(" reported error")
	} else {
		builder.WriteString("Opaque data handler reported error")
	}
	if err.PropagatedError != nil {
		builder.WriteString(": ")
		builder.WriteString(err.PropagatedError.Error())
	}
	return builder.String()
}

type UnionDiscriminantError struct {
	Discriminant uint32
	HandlerName string
}

func(err *UnionDiscriminantError) Error() string {
	var builder strings.Builder
	if len(err.HandlerName) > 0 {
		builder.WriteString(err.HandlerName)
		builder.WriteString(" reported unrecognized discriminant: ")
	} else {
		builder.WriteString("Tagged union reported unrecognized discriminant: ")
	}
	builder.WriteString(strconv.FormatUint(uint64(err.Discriminant), 10))
	return builder.String()
}
