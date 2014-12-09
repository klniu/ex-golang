package main

import (
	"github.com/favframework/debug"
)

func main() {
	a := make(map[string]int64)

	a["A"] = 1
	a["B"] = 2

	debug.Dump(a)
}
