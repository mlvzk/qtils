package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

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
	parser.AddBoolean("show-ends")
	parser.AddAliases("show-ends", "E")
	parser.AddBoolean("number")
	parser.AddAliases("number", "n")
	parser.AddBoolean("show-nonprinting")
	parser.AddAliases("show-nonprinting", "v")

	command := parser.Parse(os.Args)

	_, showEnds := command.Args["show-ends"]
	_, number := command.Args["number"]
	_, showNonprinting := command.Args["show-nonprinting"]

	var output io.Writer = os.Stdout
	if showEnds {
		output = newShowEndsProxy(output)
	}
	if number {
		output = newNumberProxy(output)
	}
	if showNonprinting {
		output = newShowNonprintingProxy(output)
	}

	for _, arg := range command.Positionals {
		var file *os.File
		var err error
		if arg == "-" {
			file = os.Stdin
		} else {
			if file, err = os.Open(arg); err != nil {
				log.Fatalf("Open file '%s' error: %v\n", arg, err)
			}
		}

		if _, err := io.Copy(output, file); err != nil {
			log.Fatalf("Copying file '%s' error: %v\n", arg, err)
		}
	}

	if len(command.Positionals) == 0 {
		if _, err := io.Copy(output, os.Stdin); err != nil {
			log.Fatalf("Copying from stdin error: %v\n", err)
		}
	}
}

type showEndsProxy struct {
	original io.Writer
}

func newShowEndsProxy(original io.Writer) showEndsProxy {
	return showEndsProxy{
		original,
	}
}

func (proxy showEndsProxy) Write(p []byte) (int, error) {
	var pos int

	buffer := bytes.Buffer{}
	buffer.Grow(len(p))
	for i, b := range p {
		if b == '\n' {
			buffer.Write(p[pos:i])
			buffer.Write([]byte("$\n"))
			pos = i + 1
		}
	}

	if pos < len(p) {
		buffer.Write(p[pos:])
		pos = len(p)
	}

	_, err := buffer.WriteTo(proxy.original)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

type numberProxy struct {
	original io.Writer
	lineNum  int
}

func newNumberProxy(original io.Writer) numberProxy {
	return numberProxy{
		original,
		1,
	}
}

func (proxy numberProxy) Write(bytes []byte) (int, error) {
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
		line := append([]byte(leftPad(strconv.Itoa(proxy.lineNum), " ", 5)+"\t"), bytes[pos:]...)
		proxy.original.Write(line)
		proxy.lineNum++
		pos = len(bytes)
	}

	return pos, nil
}

type showNonprintingProxy struct {
	original io.Writer
}

func newShowNonprintingProxy(original io.Writer) showNonprintingProxy {
	return showNonprintingProxy{
		original,
	}
}

func (proxy showNonprintingProxy) Write(p []byte) (int, error) {
	buffer := bytes.Buffer{}
	buffer.Grow(len(p))
	for _, b := range p {
		if b >= utf8.RuneSelf { // not ascii
			buffer.Write([]byte("M-"))
			b = b & 0x7f // toascii
		}

		if unicode.IsControl(rune(b)) {
			buffer.WriteByte('^')
			if b == '\177' {
				b = '?'
			} else {
				b = b | 0100
			}
		}

		buffer.WriteByte(b)
	}

	_, err := buffer.WriteTo(proxy.original)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
