package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var files map[string]map[string]*os.File = make(map[string]map[string]*os.File)

func parseDft(p string, data Data, t *template.Template) {
	os.MkdirAll(filepath.Dir(p), os.ModeDir)
	f, err := os.Create(p)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = t.Execute(f, data) // os.Stdout
	if err != nil {
		log.Fatal(err)
	}
}

func parseStr(data Data, t *template.Template) string {
	var b bytes.Buffer // A Buffer needs no initialization.

	err := t.Execute(&b, data) // os.Stdout
	if err != nil {
		log.Fatal(err)
	}

	return b.String()
}

func parseMdl(p string, data Data, t *template.Template, node string, header string) {
	var err error

	if _, ok := files[node]; !ok {
		files[node] = make(map[string]*os.File)
	}

	if _, ok := files[node][data.Const["module"]]; !ok {
		log.Printf("%#v\n", filepath.Dir(p))
		err = os.MkdirAll(filepath.Dir(p), os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
		files[node][data.Const["module"]], err = os.Create(p)
		if err != nil {
			log.Fatal(err)
		}
		_, err = files[node][data.Const["module"]].WriteString(header)
		if err != nil {
			log.Fatal(err)
		}
		//defer files[node][data.Const["module"]].Close()
	}

	_, err = files[node][data.Const["module"]].WriteString(parseStr(data, t))
	if err != nil {
		log.Fatal(err)
	}
}

func closeFiles() {
	for k, _ := range files {
		for i, _ := range files[k] {
			files[k][i].Close()
		}
	}
}
