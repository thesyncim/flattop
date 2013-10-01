package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	Reset      = "\x1b[0m"
	Bright     = "\x1b[1m"
	Dim        = "\x1b[2m"
	Underscore = "\x1b[4m"
	Blink      = "\x1b[5m"
	Reverse    = "\x1b[7m"
	Hidden     = "\x1b[8m"

	FgBlack   = "\x1b[30m"
	FgRed     = "\x1b[31m"
	FgGreen   = "\x1b[32m"
	FgYellow  = "\x1b[33m"
	FgBlue    = "\x1b[34m"
	FgMagenta = "\x1b[35m"
	FgCyan    = "\x1b[36m"
	FgWhite   = "\x1b[37m"

	BgBlack   = "\x1b[40m"
	BgRed     = "\x1b[41m"
	BgGreen   = "\x1b[42m"
	BgYellow  = "\x1b[43m"
	BgBlue    = "\x1b[44m"
	BgMagenta = "\x1b[45m"
	BgCyan    = "\x1b[46m"
	BgWhite   = "\x1b[47m"
)

type cmd []string

func (cmd *cmd) add(command ...string) {
	*cmd = append(*cmd, command...)
}

func (cmd *cmd) String() string {
	return strings.Join(*cmd, " ")
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Pretty-prints the given byte slice with indention of two spaces
func Pretty(src []byte) ([]byte, error) {
	return PrettyIndent(src, "  ")
}

// Pretty-prints the given byte slice using the provided indention string
func PrettyIndent(src []byte, indent string) ([]byte, error) {
	dst := new(bytes.Buffer)
	err := json.Indent(dst, src, "", indent)
	return dst.Bytes(), err
}

func info(message ...interface{}) {
	fmt.Print(FgGreen)
	fmt.Print(message...)
	fmt.Print(Reset + "\n")

}

func exit(message ...interface{}) {
	fmt.Print(FgRed)
	fmt.Print(message...)
	fmt.Print(Reset + "\n")
	os.Exit(-1)
}
