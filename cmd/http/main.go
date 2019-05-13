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
)

func main() {
	parser := commandparser.New()
	helper := commandhelper.New()

	helper.SetName("http")
	helper.SetVersion("v0.1.0")
	helper.AddAuthors("mlvzk")

	parser.AddOption(helper.EatOption(
		commandhelper.
			NewOption("port").
			Alias("p").
			// Default("3434").
			Required().
			Description("Port which the file server will listen on").
			Build(),
		commandhelper.
			NewOption("header").
			Alias("H").
			Arrayed().
			Description("Header that will be sent with the request. Can be multiple: -H 'Content-Type: application/json' -H 'Something: 1'").
			Build(),
		commandhelper.
			NewOption("verbose").
			Alias("v").
			Arrayed().
			Boolean().
			Description("Verbosity level. Can be multiple. Most verbose: -v -v -v").
			Build(),
		commandhelper.
			NewOption("boolean").
			Alias("b").
			Boolean().
			Build(),
	)...)

	command := parser.Parse(os.Args)
	command.Args = helper.FillDefaults(command.Args)
	errors := helper.VerifyArgs(command.Args)
	for _, err := range errors {
		log.Fatalf("%v\n", err)
	}

	if command.Booleans["boolean"] {
		println("boolean")
	}

	verbosityLevel := len(command.Arrayed["verbose"])
	fmt.Println("verbosityLeveL: ", verbosityLevel)

	port, givenPort := command.Args["port"]
	if !givenPort {
		port = "3434"
	}

	if _, help := command.Args["help"]; help {
		log.Println("test")
	}

	headers := http.Header{}
	userHeaders, givenHeaders := command.Arrayed["header"]
	if !givenHeaders {
		userHeaders = []string{}
	}
	for _, userHeader := range userHeaders { // format: 'Key: Value' or 'Key:Value'
		parts := strings.Split(userHeader, ":")
		if len(parts) < 2 {
			log.Fatalf("Invalid header '%v'; Possibly missing ':' (colon)", userHeader)
		}

		value := userHeader[len(parts[0])+1:]
		if len(value) >= 1 && value[0] == ' ' {
			value = value[1:]
		}

		headers.Add(parts[0], value)
	}

	if len(command.Positionals) == 0 { // listen mode
		server := http.FileServer(http.Dir("."))
		http.Handle("/", server)

		log.Println("Listening on port " + port)
		http.ListenAndServe(":"+port, nil)
	} else { // request mode
		method := "GET"
		url := ""
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
			var err error

			filePath := command.Positionals[2]
			if filePath == "-" {
				reader = os.Stdin
			} else {
				reader, err = os.Open(filePath)
				if err != nil {
					log.Fatalf("Error opening the file %v ; Error: %v", filePath, err)
				}
			}
		}

		if len(url) < 4 || url[0:4] != "http" {
			url = "http://" + url
		}

		request, err := http.NewRequest(method, url, reader)
		if err != nil {
			log.Fatalf("Error on creating request object with %v %v ; Error: %v\n", method, url, err)
		}
		request.Header = headers

		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalf("Error on HTTP %v %v ; Error: %v\n", method, url, err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			ioutil.ReadAll(res.Body)
		}

		if _, err := io.Copy(os.Stdout, res.Body); err != nil {
			log.Fatalf("Error on reading body from HTTP GET %v ; Error: %v\n", command.Positionals[0], err)
		}
	}
}
