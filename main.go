package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var h hash.Hash
var hashFunc = "sha384"
var resUrl string

func main() {

	args := os.Args

	if len(args) == 1 {
		println("No arguments supplied. See \"sri -h\" for help.")
		os.Exit(1)
	}

	// help page
	if args[1] == "--help" || args[1] == "-h" {
		println("Usage:")
		println("  sri [OPTION] <url>")
		println("")
		println("Options:")
		println("  -sha256	Use sha256 as hash function")
		println("  -sha384	Use sha384 as hash function (default)")
		println("  -sha512	Use sha512 as hash function")
		println("")
		println("Examples:")
		println("  sri https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js")
		println("  sri -sha512 https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css")
		println("")
		os.Exit(0)
	}

	// check arguments
	if len(args) == 3 {
		resUrl = strings.ToLower(args[2])
	} else if len(args) == 2 {
		resUrl = strings.ToLower(args[1])
	} else {
		println("Invalid amount of arguments.")
		os.Exit(1)
	}

	// get defined hash function
	// default is sha384
	switch args[1] {
	case "-sha256":
		h = sha256.New()
		hashFunc = "sha256"

	case "-sha384":
		h = sha512.New384()

	case "-sha512":
		h = sha512.New()
		hashFunc = "sha512"

	default:

		// exit if the defined hash function is none of the above
		if strings.HasPrefix(args[1], "-") {
			println("Invalid hash function \"" + args[1] + "\". Available hash functions are sha256, sha384 and sha512.")
			os.Exit(1)
		}

		h = sha512.New384()
	}

	// retrieve the content of the specified url and write it into variable h
	retrieveContent(resUrl)

	// create the hash and encode it with base64
	hash64 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if strings.HasSuffix(resUrl, ".css") {
		println("<link rel=\"stylesheet\" href=\"" + resUrl + "\" integrity=\"" + hashFunc + "-" + hash64 + "\" crossorigin=\"anonymous\">")
	} else {
		println("<script src=\"" + resUrl + "\" integrity=\"" + hashFunc + "-" + hash64 + "\" crossorigin=\"anonymous\"></script>")
	}
}

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	} else {
		return true
	}
}

func exitWithError(err error) {
	println("An error occurred:")
	println(err)
	os.Exit(1)
}

func retrieveContent(url string) {

	// check url for validity
	if !isValidUrl(url) {
		println("Invalid URL given.")
		os.Exit(1)
	}

	// initialize the http client
	client := &http.Client{
		Timeout: time.Second * 10, // define http client timeout
	}

	// create the http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		exitWithError(err)
	}

	// set additional headers
	req.Header.Set("User-Agent", "SRI Tool/1.0")

	// execute the http request
	resp, err := client.Do(req)
	if err != nil {
		exitWithError(err)
	}

	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != 200 {
		println(url + ": " + resp.Status)
		os.Exit(1)
	}

	_, err = io.Copy(h, resp.Body)
	if err != nil {
		exitWithError(err)
	}
}
