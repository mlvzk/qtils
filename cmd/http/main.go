package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mlvzk/qtils/commandparser"
)

func main() {
	parser := commandparser.New()
	parser.AddAliases("port", "p")
	parser.AddArrayed("header")
	parser.AddAliases("header", "h", "H")

	command := parser.Parse(os.Args)

	port, givenPort := command.Args["port"]
	if !givenPort {
		port = "3434"
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
