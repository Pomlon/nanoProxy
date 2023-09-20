package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func main() {

	// define origin server URL
	//originServerURL, err := url.Parse("https://sqlite.org/")
	//if err != nil {
	//	log.Fatal("invalid origin server URL")
	//}

	var cnf = koanf.New(".")
	if err := cnf.Load(file.Provider("./config.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("Error loading config.yaml: %v", err)
	}

	var endpointsMap = cnf.StringMap("endpoints")
	var listenInt = cnf.String("listenInterface")
	var listenPort = cnf.String("listenPort")
	var listenURL = listenInt + ":" + listenPort
	fmt.Println("Proxy listening on", listenURL)

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[reverse proxy server] received request at: %s\n", time.Now())

		var path = strings.ReplaceAll(req.URL.Path, "/", "")
		//fmt.Println(path)
		var ruleMatch = false
		for key, value := range endpointsMap {
			if key == path {
				var ruleDest, err = url.Parse(value)
				if err != nil {
					fmt.Println("Error while parsing url", err)
				}
				fmt.Println("Found match in rule list")
				req.Host = ruleDest.Host
				req.URL.Host = ruleDest.Host
				req.URL.Scheme = ruleDest.Scheme
				req.URL.Path = ruleDest.Path

				ruleMatch = true
				break
			}
		}

		if !ruleMatch {
			io.Copy(rw, strings.NewReader("No matching rule found"))
			return
		}

		// set req Host, URL and Request URI to forward a request to the origin server
		//req.Host = originServerURL.Host
		//req.URL.Host = originServerURL.Host
		//req.URL.Scheme = originServerURL.Scheme
		req.RequestURI = ""

		fmt.Println("Fetching ", req.URL)

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		// save the response from the origin server
		originServerResponse, err := http.DefaultClient.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}

		// return response to the client

		io.Copy(rw, originServerResponse.Body)
	})

	log.Fatal(http.ListenAndServe(listenURL, reverseProxy))
}
