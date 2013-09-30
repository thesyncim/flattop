package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type cmd []string

func (cmd *cmd) add(command ...string) {
	*cmd = append(*cmd, command...)
}

func (cmd *cmd) String() string {
	return strings.Join(*cmd, " ")
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
	fmt.Println(message...)
}
