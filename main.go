package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Rules struct {
	XMLName xml.Name `xml:"rules"`
	Rules   []Rule   `xml:"rule"`
}

type Rule struct {
	XMLName xml.Name `xml:"rule"`
	Adds    []Add    `xml:"conditions>add"`
	Action  Action   `xml:"action"`
}

func (r Rule) String() string {
	return fmt.Sprintf("Add - %v, Action - %v", r.Adds, r.Action)
}

type Action struct {
	Url string `xml:"url,attr"`
}

func (ac Action) String() string {
	return fmt.Sprintf("Expected: %s", ac.Url)
}

type Add struct {
	XMLName xml.Name `xml:"add"`
	Pattern string   `xml:"pattern,attr"`
}

func (a Add) String() string {
	return fmt.Sprintf("Pattern: %s", a.Pattern)
}

func main() {
	f, err := os.Open("rulesXml.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, _ := ioutil.ReadAll(f)

	var r Rules
	if err := xml.Unmarshal(data, &r); err != nil {
		panic(err)
	}

	for i := 0; i < len(r.Rules); i++ {
		fmt.Println(r.Rules[i])
	}
}
