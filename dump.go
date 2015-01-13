package main

import (
	"github.com/liudng/godump"
)

func main() {
	a := make(map[string]int64)

	a["A"] = 1
	a["B"] = 2

	godump.Dump(a)
}
