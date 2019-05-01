package main

import (
	"io"
	"log"
	"os"

	"github.com/mlvzk/qtils/commandparser"
)

func main() {
	parser := commandparser.New()
	parser.AddBoolean("show-ends")

	command := parser.Parse(os.Args)

	_, showEnds := command.Args["show-ends"]

	var output io.Writer = os.Stdout
	if showEnds {
		output = NewShowEndsProxy(output)
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

	return pos, nil
}
