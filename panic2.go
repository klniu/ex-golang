package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 10; i++ {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered: ", r)
			}
		}()
		if i == 3 || i == 7 {
			panic("i = 3")
		}
	}
}
