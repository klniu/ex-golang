package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	a := 1
	b := a + 2
	b = b + 1
	fmt.Println(a)
	fmt.Println(b)

	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(filepath.Dir(pathAbs))
}
