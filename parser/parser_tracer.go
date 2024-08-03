package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 1

const traceIdent = "\t"

func IdentLevel() string {
	return strings.Repeat(traceIdent, traceLevel)
}

func TracePrint(s string) {
	fmt.Printf("%s%s\n", IdentLevel(), s)
}

func IncIdent() { traceLevel++ }
func DecIdent() { traceLevel-- }
func Trace(msg string) string {
	IncIdent()
	TracePrint("Begin " + msg)
	return msg
}

func UnTrace(msg string) {
	TracePrint("END " + msg)
	DecIdent()
}
