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

	tBiz := template.Must(template.New("biz").Funcs(funcMap).Parse(tplBiz))
	tCtr := template.Must(template.New("ctr").Funcs(funcMap).Parse(tplCtr))
	tEnt := template.Must(template.New("ent").Funcs(funcMap).Parse(tplEnt))
	tIom := template.Must(template.New("iom").Funcs(funcMap).Parse(tplIom))
	tMod := template.Must(template.New("mod").Funcs(funcMap).Parse(tplMod))
	tTab := template.Must(template.New("tab").Funcs(funcMap).Parse(tplTab))

	tAdd := template.Must(template.New("add").Funcs(funcMap).Parse(tplAdd))
	tBrowse := template.Must(template.New("browse").Funcs(funcMap).Parse(tplBrowse))
	tDetail := template.Must(template.New("detail").Funcs(funcMap).Parse(tplDetail))
	tEdit := template.Must(template.New("edit").Funcs(funcMap).Parse(tplEdit))

	for _, table := range tables {
		data := allColumns(table)

		// biz
		p := console.WorkingDir + "/backend/biz/" + table.TableName.String + ".go"
		parseTpl(data, p, tBiz)

		// ctr
		p = console.WorkingDir + "/backend/ctr/" + table.TableName.String + ".go"
		parseTpl(data, p, tCtr)

		// ent
		p = console.WorkingDir + "/backend/ent/" + table.TableName.String + ".go"
		parseTpl(data, p, tEnt)

		// iom
		p = console.WorkingDir + "/backend/iom/" + table.TableName.String + ".go"
		parseTpl(data, p, tIom)

		// mod
		p = console.WorkingDir + "/backend/mod/" + table.TableName.String + ".go"
		parseTpl(data, p, tMod)

		// tab
		p = console.WorkingDir + "/backend/tab/" + table.TableName.String + ".go"
		parseTpl(data, p, tTab)

		// add
		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_add.jsx"
		parseTpl(data, p, tAdd)

		// browse
		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_browse.jsx"
		parseTpl(data, p, tBrowse)

		// detail
		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_detail.jsx"
		parseTpl(data, p, tDetail)

		// edit
		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_edit.jsx"
		parseTpl(data, p, tEdit)
	}
}
