package main

import (
	"github.com/zhgo/console"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	log.SetFlags(log.Lshortfile)

	tables := allTables()

	tController := template.Must(template.New("controller").Funcs(funcMap).Parse(tplController))
	tModel := template.Must(template.New("model").Funcs(funcMap).Parse(tplModel))
	tBrowse := template.Must(template.New("browse").Funcs(funcMap).Parse(tplBrowse))
	tDetail := template.Must(template.New("detail").Funcs(funcMap).Parse(tplDetail))
	tAdd := template.Must(template.New("add").Funcs(funcMap).Parse(tplAdd))
	tEdit := template.Must(template.New("edit").Funcs(funcMap).Parse(tplEdit))

	for _, table := range tables {
		data := allColumns(table)

		p := console.WorkingDir + "/backend/" + data.Const["module"] + "/controller/" + data.Const["entity"] + ".go"
		parseTpl(data, p, tController)

		p = console.WorkingDir + "/backend/" + data.Const["module"] + "/" + data.Const["entity"] + ".go"
		parseTpl(data, p, tModel)

		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_browse.jsx"
		parseTpl(data, p, tBrowse)

		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_detail.jsx"
		parseTpl(data, p, tDetail)

		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_add.jsx"
		parseTpl(data, p, tAdd)

		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_edit.jsx"
		parseTpl(data, p, tEdit)
	}
}
