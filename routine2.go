package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	sem := make(chan int, 10)

	for i := 0; i < 10; i++ {
		go goFunc(i, sem)
	}

	for i := 0; i < 10; i++ {
		<-sem
	}
}

func goFunc(s int, sem chan int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("goFunc: %d %d\n", s, i)
	}
	sem <- 0
}
