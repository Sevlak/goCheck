package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

var timeout int
var file string
var wg sync.WaitGroup

func init() {
	flag.IntVar(&timeout, "timeout", 60, "sets the timeout for the request")
	flag.StringVar(&file, "filename", "", "sspecifies the .xml file to be checked")
}

func main() {
	flag.Parse()

	if file == "" {
		fmt.Fprint(os.Stderr, "Missing filename flag. Use -filename [filename].")
		os.Exit(2)
	}

	r := getXMLRules(os.Args[1])
	status := make(chan string, countAdds(r))
	c := Client{&http.Client{Timeout: time.Duration(timeout) * time.Second}}

	for _, rule := range r.Rules {
		for _, add := range rule.Adds {
			wg.Add(1)
			go c.checkUrl(add.Pattern, rule.Action.Url, status, &wg)
		}
	}

	wg.Wait()
	close(status)

	for rl := range status {
		fmt.Println(rl)
	}
}

//Count how many adds are in total
func countAdds(r Rules) (x int) {
	for _, rule := range r.Rules {
		for range rule.Adds {
			x++
		}
	}

	return x
}

//Unmarshal the .xml file with the redirect rules
func getXMLRules(filepath string) Rules {
	f, err := os.Open("rulesXml.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var r Rules
	if err := xml.Unmarshal(data, &r); err != nil {
		panic(err)
	}

	return r
}
