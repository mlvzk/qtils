package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mlvzk/qtils/commandparser"
)

func leftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}
func rightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}

func main() {
	parser := commandparser.New()
	parser.AddBoolean("show-ends", "number")

	command := parser.Parse(os.Args)

	_, showEnds := command.Args["show-ends"]
	_, number := command.Args["number"]

	var output io.Writer = os.Stdout
	if showEnds {
		output = NewShowEndsProxy(output)
	}
	if number {
		output = NewNumberProxy(output)
	}

	for _, arg := range command.Positionals {
		file, err := os.Open(arg)
		if err != nil {
			log.Fatalf("Open file '%s' error: %v\n", arg, err)
		}

		if _, err = io.Copy(output, file); err != nil {
			log.Fatalf("Copying file '%s' error: %v\n", arg, err)
		}
	}
}

type ShowEndsProxy struct {
	original io.Writer
}

func NewShowEndsProxy(original io.Writer) ShowEndsProxy {
	return ShowEndsProxy{
		original,
	}
}

func (proxy ShowEndsProxy) Write(bytes []byte) (int, error) {
	var pos int

	for i, b := range bytes {
		if b == '\n' {
			proxy.original.Write(bytes[pos:i])
			proxy.original.Write([]byte("$\n"))
			pos = i + 1
		}
	}

	if pos < len(bytes) {
		proxy.original.Write(bytes[pos:])
		pos = len(bytes)
	}

	return pos, nil
}

type NumberProxy struct {
	original io.Writer
	lineNum  int
}

func NewNumberProxy(original io.Writer) NumberProxy {
	return NumberProxy{
		original,
		1,
	}
}

func (proxy NumberProxy) Write(bytes []byte) (int, error) {
	var pos int

	for i, b := range bytes {
		if b == '\n' {
			line := append([]byte(leftPad(strconv.Itoa(proxy.lineNum), " ", 5)+"\t"), bytes[pos:i+1]...)
			proxy.original.Write(line)
			proxy.lineNum++
			pos = i + 1
		}
	}

	if pos < len(bytes) {
		proxy.original.Write(bytes[pos:])
		pos = len(bytes)
	}

	return pos, nil
}
