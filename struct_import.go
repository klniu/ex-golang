package main

import (
	"fmt"
)

type a struct {

}

type b struct {
	a
}

func (c *a)Say() {
	fmt.Println("a.Say")
}

func (d *b)Say() {
	fmt.Println("b.Say")
}

func main() {
	e := b{}
	e.Say()
	e.a.Say()
}