package main

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

const tplEntHeader = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package ent

import (
    "time"
)

`

const tplEnt = `type {{ .Const.Table }} struct { {{ range $key, $column := .Columns }}
    {{ UnderscoreToCamelcase $column.ColumnName.String }}    {{ getType $column.DataType.String }}    {{ $.Const.Backquote }}json:"{{ $column.ColumnName.String }}"{{ if eq $column.ColumnKey.String "PRI" }} pk:"true"{{ end }}{{ $.Const.Backquote }}{{ end }}
}

`

const tplIomHeader = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package iom

import (
    "recom/backend/mod"
)

`

const tplIom = `var {{ .Const.Table }} = mod.New{{ .Const.Table }}()

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

const tplTabHeader = `// Copyright 2015 The recom Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package tab

import (
    "github.com/zhgo/db"
    "recom/backend/ent"
)

`

const tplTab = `var {{ .Const.Table }} = db.NewTable("{{ .Table.TableName.String }}", new(ent.{{ .Const.Table }}))

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
