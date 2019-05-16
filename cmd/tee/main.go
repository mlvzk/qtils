package main

import (
	"io"
	"log"
	"os"

	"github.com/mlvzk/qtils/commandparser"
	"github.com/mlvzk/qtils/commandparser/commandhelper"
	"github.com/pkg/errors"
)

func main() {
	parser := commandparser.New()
	helper := commandhelper.New()

	helper.SetName("tee")
	helper.SetVersion("alpha")
	helper.AddAuthor("mlvzk")

	parser.AddOption(helper.EatOption(
		commandhelper.
			NewOption("append").
			Alias("a").
			Boolean().
			Build(),
	)...)

	command, _ := parser.Parse(os.Args)

	shouldAppend := command.Booleans["append"]

	var err error
	files := make([]io.Writer, len(command.Positionals)+1)
	files[len(command.Positionals)] = os.Stdout
	for i, arg := range command.Positionals {
		if shouldAppend {
			if files[i], err = os.OpenFile(arg, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
				log.Fatalln(errors.Wrapf(err, "failed to open file '%s'", arg))
			}
			continue
		}

		if files[i], err = os.Create(arg); err != nil {
			log.Fatalln(errors.Wrapf(err, "failed to open file '%s'", arg))
		}
	}

	writer := io.MultiWriter(files...)
	_, err = io.Copy(writer, os.Stdin)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to copy to files from standard input"))
	}

	for _, file := range files {
		file.(io.Closer).Close()
	}
}
