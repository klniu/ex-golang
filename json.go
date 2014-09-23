package main

import (
	"strings"
	"io"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

const jsonStream = `
    {"Name": "Ed", "Text": "Knock knock."}
    {"Name": "Sam", "Text": "Who's there?"}
    {"Name": "Ed", "Text": "Go fmt."}
    {"Name": "Sam", "Text": "Go fmt who?"}
    {"Name": "Ed", "Text": "Go fmt yourself!"}
`

type JsonMessage struct {
	Name, Text string
}

func jsonString(txt string){
	dec := json.NewDecoder(strings.NewReader(txt))
	for {
		var m JsonMessage
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", m.Name, m.Text)
	}
}

func jsonFile() {
	txt, err := ioutil.ReadFile("examples\\json.txt")
	if err != nil {
		log.Fatal(err)
	}

	jsonString(string(txt))
}

func main() {
	jsonString(jsonStream)
	jsonFile()
}
