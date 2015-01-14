package main

import (
	"bytes"
	"fmt"
)

func main() {
	a := []byte(`11<2div id="title" class="flol xiaoqu f16"></div>3`)
	b := []byte(`<div id="title" class="flol xiaoqu f16"></div>`)

	c := bytes.Index(a, b)

	fmt.Println(c)
}
