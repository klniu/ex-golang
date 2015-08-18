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
	tModHeader := template.Must(template.New("modHeader").Funcs(funcMap).Parse(tplModHeader))
	tMod1 := template.Must(template.New("mod1").Funcs(funcMap).Parse(tplMod1))
	tMod2 := template.Must(template.New("mod2").Funcs(funcMap).Parse(tplMod2))
	tMod3 := template.Must(template.New("mod3").Funcs(funcMap).Parse(tplMod3))

	tAdd := template.Must(template.New("add").Funcs(funcMap).Parse(tplAdd))
	tBrowse := template.Must(template.New("browse").Funcs(funcMap).Parse(tplBrowse))
	tDetail := template.Must(template.New("detail").Funcs(funcMap).Parse(tplDetail))
	tEdit := template.Must(template.New("edit").Funcs(funcMap).Parse(tplEdit))

	mods := make(map[string][4]string)

	for _, table := range tables {
		data := allColumns(table)

		// biz
		p := console.WorkingDir + "/backend/biz/" + data.Const["module"] + ".go"
		parseMdl(p, data, tBiz, "biz", tplBizHeader)

		// ctr
		p = console.WorkingDir + "/backend/ctr/" + data.Const["table"] + ".go"
		parseDft(p, data, tCtr)

		// ent
		p = console.WorkingDir + "/backend/ent/" + data.Const["module"] + ".go"
		parseMdl(p, data, tEnt, "ent", tplEntHeader)

		// mod
		m := mods[data.Const["module"]]
		m[0] = parseStr(data, tModHeader)
		m[1] += parseStr(data, tMod1)
		m[2] += parseStr(data, tMod2)
		m[3] += parseStr(data, tMod3)
		mods[data.Const["module"]] = m

		// add
		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_add.jsx"
		parseDft(p, data, tAdd)

		// browse
		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_browse.jsx"
		parseDft(p, data, tBrowse)

		// detail
		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_detail.jsx"
		parseDft(p, data, tDetail)

		// edit
		p = console.WorkingDir + "/frontend/" + data.Const["module"] + "/" + data.Const["entity"] + "_edit.jsx"
		parseDft(p, data, tEdit)
	}

	// mod
	for k, v := range mods {
		p := console.WorkingDir + "/backend/mod/" + k + ".go"

		err := os.MkdirAll(filepath.Dir(p), os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(p)
		if err != nil {
			log.Fatal(err)
		}

		for _, str := range v {
			_, err := f.WriteString(str)
			if err != nil {
				log.Fatal(err)
			}
		}

		f.Close()
	}

	closeFiles()
}
