package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

var TIMEOUT time.Duration = 30 * time.Second
var wg sync.WaitGroup

func main() {
	r := getXMLRules(os.Args[1])
	status := make(chan string, 5000)
	c := Client{&http.Client{Timeout: TIMEOUT}}

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
