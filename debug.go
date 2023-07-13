package goxdr

import (
	"os"
	"fmt"
)

const debugPrefix string = "***DEBUG[github.com/UncleSniper/goxdr]: "

var debugOn bool = len(os.Getenv("GO_USDO_XDR_DEBUG")) > 0

func debugf(format string, args ...any) {
	if debugOn {
		fmt.Fprint(os.Stderr, debugPrefix)
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func debugln(args ...any) {
	if debugOn {
		fmt.Fprint(os.Stderr, debugPrefix)
		fmt.Fprintln(os.Stderr, args...)
	}
}
