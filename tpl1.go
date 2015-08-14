package main

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
