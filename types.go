package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sync"
)

type Client struct {
	*http.Client
}

type Link struct {
	Pattern  string
	Status   string
	Found    string
	Expected string
}

func (c *Client) checkUrl(pattern, action string, status chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := c.Get("https://" + pattern)
	if err != nil {
		status <- fmt.Sprintf("%s,%s,%s,%s", pattern, "TIMEOUT", "TIMEOUT", action)
		return
	}
	defer resp.Body.Close()

	found := resp.Request.URL.String()
	status <- fmt.Sprintf("%s,%s,%s,%s", pattern, resp.Status, found, action)
}

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
