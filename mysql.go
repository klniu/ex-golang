package main

import (
	"database/sql"
	"github.com/zhgo/console"
	"github.com/zhgo/db"
	"log"
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

var server *db.Server = db.NewServer("mysql-1", "mysql", "root@tcp(127.0.0.1:3306)/information_schema?charset=utf8")

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

// Create a FuncMap with which to register the function.
var funcMap template.FuncMap = template.FuncMap{
	"UnderscoreToCamelcase": console.UnderscoreToCamelcase,
	"CamelcaseToUnderscore": console.CamelcaseToUnderscore,
	"getType":               getType,
}

func getType(str string) string {
	v, ok := typeMap[str]
	if !ok {
		log.Fatalf("typeMap[%s] not exists.", str)
	}
	return v
}

func allTables() []Table {
	db.Env = 3
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
