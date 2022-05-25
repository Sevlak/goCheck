package main

import (
	"encoding/xml"
	"fmt"
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
