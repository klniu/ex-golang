package main

import (
	"fmt"
)

type A struct {
	ID int64
}

type B struct {
	A    // Import A
	ID   int64
	Name string
}

func (a *A) Init() {
	a.ID = 128
}

func (b *B) Init() {
	b.A.Init()
	b.ID = 256
}

func (a *A) Say() {
	fmt.Println("A.Say: ", a.ID)
}

func (b *B) Say() {
	b.A.Say()
	fmt.Println("B.Say: ", b.ID, b.A.ID)
}

func main() {
	b := B{}
	b.Init()
	b.Say()
	//b.A.Say()
}
