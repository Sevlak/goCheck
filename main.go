package main

import (
	"bufio"
	"encoding/csv"
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
var filter bool
var wg sync.WaitGroup

func init() {
	flag.IntVar(&timeout, "timeout", 60, "sets the timeout for the request in seconds")
	flag.StringVar(&file, "filename", "", "specifies the .xml file to be checked")
	flag.BoolVar(&filter, "filter", false, "filter results where the found url doesn't match the expected url")
}

func main() {
	flag.Parse()

	if file == "" {
		fmt.Fprint(os.Stderr, "Missing filename flag. Use -filename [filename].\n")
		os.Exit(2)
	}

	r := getXMLRules(file)
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

	writeToCsv(status)

	if filter {
		createFilteredCsv()
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
	f, err := os.Open(filepath)
	check(err)
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	check(err)

	var r Rules
	if err := xml.Unmarshal(data, &r); err != nil {
		panic(err)
	}

	return r
}

func writeToCsv(status chan string) {
	f, err := os.Create("results.csv")
	defer f.Close()

	check(err)

	w := bufio.NewWriter(f)
	w.WriteString("pattern,status,urlfound,urlexpected\n") //csv headers
	for s := range status {
		w.WriteString(s + "\n")
	}

	w.Flush()
}

func filterWrongRedirects(filepath string) []Link {
	f, err := os.Open(filepath)
	check(err)

	r := csv.NewReader(f)
	r.FieldsPerRecord = 4

	results, err := r.ReadAll()
	check(err)

	f.Close()

	var wrong []Link
	for _, line := range results {
		if line[2] != line[3] && line[1] != "404 Status Not Found" {
			wrong = append(wrong, Link{Pattern: line[0], Status: line[1],
				Found: line[2], Expected: line[3]})
		}
	}

	return wrong
}

func createFilteredCsv() {
	wrong := filterWrongRedirects("results.csv")
	f, err := os.Create("filtered_results.csv")
	check(err)

	w := csv.NewWriter(f)
	for _, record := range wrong {
		err = w.Write([]string{record.Pattern, record.Status, record.Found, record.Expected})
		check(err)
	}

	w.Flush()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
