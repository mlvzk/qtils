package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mlvzk/qtils/commandparser"
	"github.com/mlvzk/qtils/commandparser/commandhelper"
	"github.com/pkg/errors"
)

func main() {
	parser := commandparser.New()
	helper := commandhelper.New()

	helper.SetName("http")
	helper.SetVersion("alpha")
	helper.AddAuthor("mlvzk")

	parser.AddOption(helper.EatOption(
		commandhelper.NewOption("help").Alias("h").Boolean().Description("Prints this page"),
		commandhelper.
			NewOption("protocol").
			Description("protocol").
			Required().
			ValidateBind(commandhelper.ValidateSelection("http", "https", "gopher")),
		// TODO: remove this option
		commandhelper.
			NewOption("port").
			Alias("p").
			Default("3434").
			Required().
			ValidateBind(commandhelper.ValidateInt).
			Description("Port which the file server will listen on"),
		commandhelper.
			NewOption("header").
			Alias("H").
			Arrayed().
			Validate(func(value string) error {
				if parts := strings.Split(value, ":"); len(parts) < 2 {
					return errors.New("invalid header value, missing `:`")
				}
				return nil
			}).
			Description("Header that will be sent with the request. Can be multiple: -H 'Content-Type: application/json' -H 'Something: 1'"),
	)...)

	command, err := parser.Parse(os.Args)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to parse arguments"))
	}

	if command.Booleans["help"] {
		fmt.Print(helper.Help())
		os.Exit(1)
	}

	command.Args = helper.FillDefaults(command.Args)
	errs := helper.Verify(command.Args, command.Arrayed)
	for _, err := range errs {
		log.Println(err)
	}
	if len(errs) != 0 {
		return
	}

	port := command.Args["port"]

	headers := http.Header{}
	userHeaders := command.Arrayed["header"]
	for _, userHeader := range userHeaders { // format: 'Key: Value' or 'Key:Value'
		parts := strings.Split(userHeader, ":")
		if len(parts) < 2 {
			log.Fatalf("Invalid header '%v'; Possibly missing ':' (colon)\n", userHeader)
		}

		value := userHeader[len(parts[0])+1:]
		if len(value) >= 1 && value[0] == ' ' {
			value = value[1:]
		}

		headers.Add(parts[0], value)
	}

	if len(command.Positionals) == 0 { // listen mode
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalln(errors.Wrap(err, "failed to get working directory"))
		}

		server := http.FileServer(http.Dir(dir))
		http.Handle("/", server)

		log.Println("Listening on port " + port)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Println(errors.Wrap(err, "failed to listen and serve"))
		}
	} else { // request mode
		method := "GET"
		var url string
		var reader io.ReadCloser

		switch len(command.Positionals) {
		case 1: // http api.ipify.org
			url = command.Positionals[0]
		case 2: // http GET api.ipify.org
			method = strings.ToUpper(command.Positionals[0])
			url = command.Positionals[1]
		case 3: // http GET api.ipify.org ./body.txt
			method = strings.ToUpper(command.Positionals[0])
			url = command.Positionals[1]

			filePath := command.Positionals[2]
			if filePath == "-" {
				reader = os.Stdin
			} else {
				var err error
				reader, err = os.Open(filePath)
				if err != nil {
					log.Fatalln(errors.Wrapf(err, "failed to open file '%s'", filePath))
				}
			}
		}

		if len(url) < 4 || url[0:4] != "http" {
			url = "http://" + url
		}

		request, err := http.NewRequest(method, url, reader)
		if err != nil {
			log.Fatalln(errors.Wrapf(err, "failed to create request object with %v %v", method, url))
		}
		request.Header = headers

		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(errors.Wrapf(err, "failed to send http request with %v %v", method, url))
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			ioutil.ReadAll(res.Body)
		}

		if _, err := io.Copy(os.Stdout, res.Body); err != nil {
			log.Fatalln(errors.Wrapf(err, "failed to read body from %v %v", res.Request.Method, res.Request.URL.String()))
		}
	}
}
