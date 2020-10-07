package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var hashName string
var useSha256 bool
var useSha384 bool
var useSha512 bool
var urls = make([]string, 0)
var client = &http.Client{
	Timeout: time.Second * 10,
}
var wg = &sync.WaitGroup{}

func main() {
	flag.BoolVar(&useSha256, "sha256", false, "Use sha256 as hash function")
	flag.BoolVar(&useSha384, "sha384", true, "Use sha384 as hash function")
	flag.BoolVar(&useSha512, "sha512", false, "Use sha512 as hash function")
	flag.Usage = usage
	flag.Parse()

	// if no urls are supplied as arguments check if stdin contains urls
	if len(flag.Args()) == 0 {
		fi, err := os.Stdin.Stat()
		if err != nil {
			exitWithError(err)
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			fmt.Println("No urls supplied. See \"sri -h\" for help.")
			os.Exit(1)
		}

		scan := bufio.NewScanner(os.Stdin)
		for {
			scan.Scan()
			input := scan.Text()

			if len(input) != 0 {
				urls = append(urls, input)
			} else {
				break
			}
		}
	} else {
		urls = flag.Args()
	}

	// set hash name
	if useSha256 {
		hashName = "sha256"
	} else if useSha512 {
		hashName = "sha512"
	} else {
		hashName = "sha384"
	}

	for _, v := range urls {
		wg.Add(1)
		go generateSri(v)
	}

	// wait for all goroutines to finish
	wg.Wait()
}

func generateSri(uri string) {
	var h hash.Hash

	// initialize the chosen hash function
	if useSha256 {
		h = sha256.New()
	} else if useSha512 {
		h = sha512.New()
	} else {
		h = sha512.New384()
	}

	// retrieve the content of the specified url and write it into variable h
	d, err := retrieveContent(uri)
	if err != nil {
		exitWithError(err)
	}

	// create the hash and encode it with base64
	if _, err := h.Write(d); err != nil {
		exitWithError(err)
	}
	sriHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if strings.HasSuffix(uri, ".css") {
		fmt.Println("<link rel=\"stylesheet\" href=\"" + uri + "\" integrity=\"" + hashName + "-" + sriHash + "\" crossorigin=\"anonymous\">")
	} else {
		fmt.Println("<script src=\"" + uri + "\" integrity=\"" + hashName + "-" + sriHash + "\" crossorigin=\"anonymous\"></script>")
	}

	wg.Done()
}

func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  sri [OPTION] [<url1> <url2> ... <urlN>]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -sha256	Use sha256 as hash function")
	fmt.Println("  -sha384	Use sha384 as hash function (default)")
	fmt.Println("  -sha512	Use sha512 as hash function")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  sri https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js")
	fmt.Println("  sri -sha512 https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css")
	fmt.Println("")
}

func retrieveContent(uri string) ([]byte, error) {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return nil, err
	}

	// create the http request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	// set user agent
	req.Header.Set("User-Agent", "SRI cli/0.0.4 (https://github.com/baretools/sri)")

	// execute the http request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// check status code
	if resp.StatusCode != 200 {
		return nil, errors.New(uri + " returned status code " + strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
