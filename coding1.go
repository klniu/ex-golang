package main

import (
	"database/sql"
	"github.com/zhgo/console"
	"github.com/zhgo/db"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Table struct {
	TableName    sql.NullString `json: "TABLE_NAME"`
	TableComment sql.NullString `json: "TABLE_COMMENT"`
}

type Column struct {
	ColumnName             sql.NullString `json: "COLUMN_NAME"`
	DataType               sql.NullString `json: "DATA_TYPE"`
	NumericPrecision       sql.NullInt64  `json: "NUMERIC_PRECISION"`
	NumericScale           sql.NullInt64  `json: "NUMERIC_SCALE"`
	CharacterMaximumLength sql.NullInt64  `json: "CHARACTER_MAXIMUM_LENGTH"`
	ColumnDefault          sql.NullString `json: "COLUMN_DEFAULT"`
	ColumnType             sql.NullString `json: "COLUMN_TYPE"`
	ColumnComment          sql.NullString `json: "COLUMN_COMMENT"`
	ColumnKey              sql.NullString `json: "COLUMN_KEY"`
}

type Data struct {
	Const   map[string]string
	Table   Table
	Columns []Column
}

func allTables() []Table {
	tables := []Table{}
	q := server.Select("TABLE_NAME", "TABLE_COMMENT")
	err := q.From("TABLES").Where(q.Eq("TABLE_SCHEMA", "recom")).Rows(&tables)
	if err != nil {
		log.Fatal(err)
	}
	return tables
}

func allColumns(table Table) Data {
	columns := []Column{}
	q := server.Select("COLUMN_NAME", "DATA_TYPE", "NUMERIC_PRECISION", "NUMERIC_SCALE", "CHARACTER_MAXIMUM_LENGTH", "COLUMN_DEFAULT", "COLUMN_TYPE", "COLUMN_COMMENT", "COLUMN_KEY")
	err := q.From("COLUMNS").Where(q.Eq("TABLE_SCHEMA", "recom"), q.AndEq("TABLE_NAME", table.TableName.String)).Rows(&columns)
	if err != nil {
		log.Fatal(err)
	}

	modules := strings.SplitN(table.TableName.String, "_", 2)
	con := map[string]string{
		"module":    modules[0],
		"Module":    console.UnderscoreToCamelcase(modules[0]),
		"entity":    modules[1],
		"Entity":    console.UnderscoreToCamelcase(modules[1]),
		"Backquote": "`",
		"BracketL":  "{",
		"BracketR":  "}",
	}
	data := Data{con, table, columns}

	return data
}

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

func getType(str string) string {
	v, ok := typeMap[str]
	if !ok {
		log.Fatalf("typeMap[%s] not exists.", str)
	}
	return v
}

// Create a FuncMap with which to register the function.
var funcMap template.FuncMap = template.FuncMap{
	"UnderscoreToCamelcase": console.UnderscoreToCamelcase,
	"CamelcaseToUnderscore": console.CamelcaseToUnderscore,
	"getType":               getType,
}

var typeMap map[string]string = map[string]string{
	"bigint":    "int64",
	"char":      "string",
	"date":      "string",
	"decimal":   "float64",
	"enum":      "string",
	"int":       "int64",
	"smallint":  "int64",
	"tinyint":   "int64",
	"timestamp": "string",
	"varchar":   "string",
}

var server *db.Server = db.NewServer("mysql-1", "mysql", "root@tcp(127.0.0.1:3306)/information_schema?charset=utf8")

func main() {
	log.SetFlags(log.Lshortfile)

	db.Env = 3
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

const tplController = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package {{ .Const.module }}_controller

import (
	"github.com/zhgo/web"
)

type {{ .Const.Entity }}Controller struct {
	// import web.Controller
	web.Controller

	// import web.Crud
	web.Crud
}

func init() {
	web.NewController("{{ .Const.Module }}", new({{ .Const.Entity }}Controller))
}
`

const tplModel = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package {{ .Const.module }}

import (
	"github.com/zhgo/db"
)

// Entity struct
type {{ .Const.Entity }}Entity struct { {{ range $key, $column := .Columns }}
	{{ UnderscoreToCamelcase $column.ColumnName.String }}    {{ getType $column.DataType.String }}    {{ $.Const.Backquote }}json:"{{ $column.ColumnName.String }}"{{ if eq $column.ColumnKey.String "PRI" }} pk:"true"{{ end }}{{ $.Const.Backquote }}{{ end }}
}

// Model struct
type {{ .Const.Entity }}Model struct {
	db.Model //Import db.Model
}

// Table
var {{ .Const.Entity }}Table = db.NewTable("{{ .Table.TableName.String }}", new({{ .Const.Entity }}Entity))

// Model
var {{ .Const.Entity }} = New{{ .Const.Entity }}Model()

// New Model
func New{{ .Const.Entity }}Model() *{{ .Const.Entity }}Model {
	return &{{ .Const.Entity }}Model{Model: db.NewModel("{{ .Const.Module }}", {{ .Const.Entity }}Table)}
}

// Model methods
`

const tplBrowse = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

"use strict";

var React = require("react");
var frontify = require("../frontify.js");
var Container = require("../container.jsx");

var {{ .Const.Entity }}List = React.createClass({
  componentDidMount: function(){
    
  },
  
  render: function(){
  	var nodes = this.props.data.map(function(item){
	  return (<tr>
	    <td>#</td>{{ range $key, $column := .Columns }}
	    <td>{{ $.Const.BracketL }}item.{{ $column.ColumnName.String }}{{ $.Const.BracketR }}</td>{{ end }}
	  </tr>);
	});

	return (<tbody>
	  {nodes}
	</tbody>);
  }
});

var {{ .Const.Entity }}Browse = React.createClass({
  componentDidMount: function(){
    
  },
  
  render: function(){
    return (<Container>
	  <h2 class="sub-header">{{ .Const.Entity }}</h2>
	  <div class="table-responsive">
	    <table class="table table-striped">
	      <thead>
	        <tr>
	          <th>#</th>{{ range $key, $column := .Columns }}
	          <th>{{ $column.ColumnComment.String }}</th>{{ end }}
	        </tr>
	      </thead>
	      <{{ .Const.Entity }}List url="/{{ .Const.module }}/{{ .Const.entity }}/browse" />
	    </table>
	  </div>
    </Container>);
  }
});

React.render(<{{ .Const.Entity }}Browse />, document.body);
`

const tplDetail = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

"use strict";

var React = require("react");
var frontify = require("../frontify.js");
var Container = require("../container.jsx");

var {{ .Const.Entity }}Detail = React.createClass({
  componentDidMount: function(){
    
  },
  
  render: function(){
    return (<Container>
    
    </Container>);
  }
});

React.render(<{{ .Const.Entity }}Detail />, document.body);
`

const tplAdd = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

"use strict";

var React = require("react");
var frontify = require("../frontify.js");
var Container = require("../container.jsx");

var {{ .Const.Entity }}Add = React.createClass({
  componentDidMount: function(){
    frontify.formValidate("#form1", function(data){
      console.log(data);
    }, function(err){
      console.log(err);
    });
  },
  
  render: function(){
    return (<Container>
    <form id="form1" action="/{{ .Const.module }}/{{ .Const.entity }}/add" method="post" className="">
      
      <button className="btn btn-lg btn-primary btn-block" type="submit">Submit</button>
    </form>
    </Container>);
  }
});

React.render(<{{ .Const.Entity }}Add />, document.body);
`

const tplEdit = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

"use strict";

var React = require("react");
var frontify = require("../frontify.js");
var Container = require("../container.jsx");

var {{ .Const.Entity }}Edit = React.createClass({
  componentDidMount: function(){
    frontify.formValidate("#form1", function(data){
      console.log(data);
    }, function(err){
      console.log(err);
    });
  },
  
  render: function(){
    return (<Container>
    <form id="form1" action="/{{ .Const.module }}/{{ .Const.entity }}/edit" method="post" className="">
      
      <button className="btn btn-lg btn-primary btn-block" type="submit">Submit</button>
    </form>
    </Container>);
  }
});

React.render(<{{ .Const.Entity }}Edit />, document.body);
`
