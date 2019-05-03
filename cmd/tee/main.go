package main

import (
	"io"
	"log"
	"os"

	"github.com/mlvzk/qtils/commandparser"
)

func main() {
	parser := commandparser.New()
	parser.AddBoolean("append")
	parser.AddAliases("append", "a")

	command := parser.Parse(os.Args)

	_, append := command.Args["append"]

	var err error
	files := make([]io.Writer, len(command.Positionals)+1)
	files[len(command.Positionals)] = os.Stdout
	for i, arg := range command.Positionals {
		if append {
			if files[i], err = os.OpenFile(arg, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
				log.Fatalf("Open file '%s' error: %v\n", arg, err)
			}
			continue
		}

		if files[i], err = os.Create(arg); err != nil {
			log.Fatalf("Open file '%s' error: %v\n", arg, err)
		}
	}

	writer := io.MultiWriter(files...)
	_, err = io.Copy(writer, os.Stdin)
	if err != nil {
		log.Fatalf("Copy error: %v\n", err)
	}

	for _, file := range files {
		file.(io.Closer).Close()
	}
}
