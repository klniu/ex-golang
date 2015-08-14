package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func parseTpl(data Data, p string, t *template.Template) {
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
