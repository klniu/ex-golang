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
		"table":     table.TableName.String,
		"Table":     console.UnderscoreToCamelcase(table.TableName.String),
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
	"timestamp": "time.Time",
	"varchar":   "string",
}

var server *db.Server = db.NewServer("mysql-1", "mysql", "root@tcp(127.0.0.1:3306)/information_schema?charset=utf8")

func main() {
	log.SetFlags(log.Lshortfile)

	db.Env = 3
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

const tplBiz = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package biz

import ()

// Biz struct
type {{ .Const.Table }} struct {
    
}

// New Biz
func New{{ .Const.Table }}() *{{ .Const.Table }} {
    return &{{ .Const.Table }}{}
}

// Biz methods
`

const tplCtr = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package ctr

import (
    "github.com/zhgo/web"
)

type {{ .Const.Table }} struct {
    // import web.Controller
    web.Controller
}

func init() {
    web.NewController("{{ .Const.Module }}", new({{ .Const.Table }}))
}

// List
func (c *{{ .Const.Table }}) List() web.Result {
    d := []ent.{{ .Const.Table }}{}

    q := iom.{{ .Const.Table }}.Select()
    err := q.Parse(c.Request.Body.Cond).Rows(&d)
    if err != nil {
        return c.Fail(err)
    }

    return c.Done(d)
}
`

const tplEnt = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package ent

import (
    "time"
)

// Entity struct
type {{ .Const.Table }} struct { {{ range $key, $column := .Columns }}
    {{ UnderscoreToCamelcase $column.ColumnName.String }}    {{ getType $column.DataType.String }}    {{ $.Const.Backquote }}json:"{{ $column.ColumnName.String }}"{{ if eq $column.ColumnKey.String "PRI" }} pk:"true"{{ end }}{{ $.Const.Backquote }}{{ end }}
}
`

const tplIom = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package iom

import (
    "recom/backend/mod"
)

// Model instance
var {{ .Const.Table }} = mod.New{{ .Const.Table }}()
`

const tplMod = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package mod

import (
    "github.com/zhgo/db"
    "recom/backend/tab"
)

// Model struct
type {{ .Const.Table }} struct {
    // Import db.Model
    db.Model
}

// New Model
func New{{ .Const.Table }}() *{{ .Const.Table }} {
    return &{{ .Const.Table }}{Model: db.NewModel("{{ .Const.Module }}", tab.{{ .Const.Table }})}
}

// Model methods
`

const tplTab = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package tab

import (
    "github.com/zhgo/db"
    "recom/backend/ent"
)

// Table
var {{ .Const.Table }} = db.NewTable("{{ .Table.TableName.String }}", new(ent.{{ .Const.Table }}))
`

const tplAdd = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

"use strict";

var React = require("react");
var frontify = require("../frontify.js");
var Container = require("../container.jsx");

var {{ .Const.Table }}Add = React.createClass({
  componentDidMount: function(){
    frontify.formValidate("#form1", function(data){
      console.log(data);
    }, function(err){
      console.log(err);
    });
  },

  render: function(){
    return (<Container>
      <h2 className="sub-header">{{ .Const.Entity }}</h2>
      <div className="text-right">
        <a className="btn btn-default" href="#{{ .Const.module }}/{{ .Const.entity }}/browse" role="button">Browse</a>
        <a className="btn btn-default" href="#{{ .Const.module }}/{{ .Const.entity }}/add" role="button">Add</a>
      </div>
      <form id="form1" action="/{{ .Const.Module }}/{{ .Const.Entity }}/Add" method="post" className="">{{ range $key, $column := .Columns }}
        <div className="form-group">
          <label for="{{ $column.ColumnName.String }}">{{ $column.ColumnComment.String }}</label>
          <input type="text" name="{{ $column.ColumnName.String }}" id="{{ $column.ColumnName.String }}" value="" className="form-control" placeholder="{{ $column.ColumnComment.String }}" />
        </div>{{ end }}

        <button className="btn btn-lg btn-primary btn-block" type="submit">Submit</button>
      </form>
    </Container>);
  }
});

React.render(<{{ .Const.Table }}Add />, document.body);
`

const tplBrowse = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

"use strict";

var React = require("react");
var frontify = require("../frontify.js");
var Container = require("../container.jsx");

var {{ .Const.Table }}Search = React.createClass({
  render: function(){
    return (<div>123</div>);
  }
});

var {{ .Const.Table }}List = React.createClass({
  apiData: function(data){
    this.setState({data: data});
  },

  apiError: function(xhr, status, err){
    console.error(status, err.toString());
  },

  handleDelete: function(v) {
    console.log(v);
  },

  handleItem: function(item){
    return (<tr>
      <td>
        <a className="btn btn-default" href={"#{{ .Const.module }}/{{ .Const.entity }}/detail/"+item.user_id} role="button">Detail</a> 
        <a className="btn btn-default" href={"#{{ .Const.module }}/{{ .Const.entity }}/edit/"+item.user_id} role="button">Edit</a> 
        <a className="btn btn-default" href="#fn" onClick={this.handleDelete.bind(null, item.user_id)} role="button">Delete</a>
      </td>{{ range $key, $column := .Columns }}
        <td>{{ $.Const.BracketL }}item.{{ $column.ColumnName.String }}{{ $.Const.BracketR }}</td>{{ end }}
    </tr>);
  },

  getInitialState: function() {
    return {data: []};
  },

  componentDidMount: function(){
    frontify.apiPost(this.props.url, {}, this.apiData, this.apiError);
  },

  render: function(){
    var nodes = this.state.data.map(this.handleItem);
    return (<tbody>{nodes}</tbody>);
  }
});

var {{ .Const.Table }}Browse = React.createClass({
  componentDidMount: function(){
    
  },

  render: function(){
    return (<Container>
      <h2 className="sub-header">{{ .Const.Entity }}</h2>
      <div className="text-right">
        <a className="btn btn-default" href="#{{ .Const.module }}/{{ .Const.entity }}/add" role="button">Add</a>
      </div>
      <{{ .Const.Table }}Search />
      <div className="table-responsive">
        <table className="table table-striped">
          <thead>
            <tr>
              <th>#</th>{{ range $key, $column := .Columns }}
              <th>{{ $column.ColumnComment.String }}</th>{{ end }}
            </tr>
          </thead>
          <{{ .Const.Table }}List url="/{{ .Const.Module }}/{{ .Const.Entity }}/List" />
        </table>
      </div>
    </Container>);
  }
});

React.render(<{{ .Const.Table }}Browse />, document.body);
`

const tplDetail = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

"use strict";

var React = require("react");
var frontify = require("../frontify.js");
var Container = require("../container.jsx");

var {{ .Const.Table }}Detail = React.createClass({
  componentDidMount: function(){
    
  },

  render: function(){
    return (<Container>
      <h2 className="sub-header">{{ .Const.Entity }}</h2>
      <div className="text-right">
        <a className="btn btn-default" href="#{{ .Const.module }}/{{ .Const.entity }}/browse" role="button">Browse</a>
        <a className="btn btn-default" href="#{{ .Const.module }}/{{ .Const.entity }}/add" role="button">Add</a>
      </div>

    </Container>);
  }
});

React.render(<{{ .Const.Table }}Detail />, document.body);
`

const tplEdit = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

"use strict";

var React = require("react");
var frontify = require("../frontify.js");
var Container = require("../container.jsx");

var {{ .Const.Table }}Edit = React.createClass({
  componentDidMount: function(){
    frontify.formValidate("#form1", function(data){
      console.log(data);
    }, function(err){
      console.log(err);
    });
  },

  render: function(){
    return (<Container>
      <h2 className="sub-header">{{ .Const.Entity }}</h2>
      <div className="text-right">
        <a className="btn btn-default" href="#{{ .Const.module }}/{{ .Const.entity }}/browse" role="button">Browse</a>
        <a className="btn btn-default" href="#{{ .Const.module }}/{{ .Const.entity }}/add" role="button">Add</a>
      </div>
      <form id="form1" action="/{{ .Const.module }}/{{ .Const.entity }}/edit" method="post" className="">{{ range $key, $column := .Columns }}
        <div className="form-group">
          <label for="{{ $column.ColumnName.String }}">{{ $column.ColumnComment.String }}</label>
          <input type="text" name="{{ $column.ColumnName.String }}" id="{{ $column.ColumnName.String }}" value="" className="form-control" placeholder="{{ $column.ColumnComment.String }}" />
        </div>{{ end }}
        
        <button className="btn btn-lg btn-primary btn-block" type="submit">Submit</button>
      </form>
    </Container>);
  }
});

React.render(<{{ .Const.Table }}Edit />, document.body);
`
