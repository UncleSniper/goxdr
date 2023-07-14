package goxdr

import (
	"strings"
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
