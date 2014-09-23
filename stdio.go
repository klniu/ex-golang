package main

import (
	"fmt"
)

func main() {
	type Employee struct {
		ID int
		Name string
	}

	e := Employee{22, "姚明"}

	fmt.Printf("%v\n\n", e)
	fmt.Printf("%+v\n\n", e)
	fmt.Printf("%T\n\n", e)
	fmt.Printf("%#v\n\n", e)
}
