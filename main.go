package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

func main() {
	rules := getXMLRules(os.Args[1])
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
