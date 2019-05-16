package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/mlvzk/qtils/commandparser"
	"github.com/mlvzk/qtils/commandparser/commandhelper"
	"github.com/mlvzk/qtils/util"
	"github.com/pkg/errors"
)

func main() {
	parser := commandparser.New()
	helper := commandhelper.New()

	helper.SetName("cat")
	helper.SetVersion("alpha")
	helper.AddAuthor("mlvzk")

	parser.AddOption(helper.EatOption(
		commandhelper.
			NewOption("show-ends").
			Boolean().
			Alias("E").
			Build(),
		commandhelper.
			NewOption("number").
			Boolean().
			Alias("n").
			Build(),
		commandhelper.
			NewOption("show-nonprinting").
			Boolean().
			Alias("v").
			Build(),
		// TODO: implement cat --squeeze-blank
	)...)

	command, err := parser.Parse(os.Args)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to parse arguments"))
	}

	showEnds := command.Booleans["show-ends"]
	number := command.Booleans["number"]
	showNonprinting := command.Booleans["show-nonprinting"]

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
				log.Fatalln(errors.Wrapf(err, "failed to open file '%s'", arg))
			}
		}

		if _, err := io.Copy(output, file); err != nil {
			log.Fatalln(errors.Wrapf(err, "failed to copy file '%s' to output", arg))
		}
	}

	if len(command.Positionals) == 0 {
		if _, err := io.Copy(output, os.Stdin); err != nil {
			log.Fatalln(errors.Wrap(err, "failed to copy from standard input to output"))
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
		return 0, errors.Wrap(err, "failed to write to original writer from proxy")
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
			line := append([]byte(util.LeftPad(strconv.Itoa(proxy.lineNum), " ", 5)+"\t"), bytes[pos:i+1]...)
			proxy.original.Write(line)
			proxy.lineNum++
			pos = i + 1
		}
	}

	if pos < len(bytes) {
		line := append([]byte(util.LeftPad(strconv.Itoa(proxy.lineNum), " ", 5)+"\t"), bytes[pos:]...)
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
	var buffer bytes.Buffer

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
		return 0, errors.Wrap(err, "failed to write to original writer from proxy")
	}

	return len(p), nil
}
