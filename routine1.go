package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 10; i++ {
		go goFunc(i)
	}
}

func goFunc(s int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("goFunc: %d %d\n", s, i)
	}
}
