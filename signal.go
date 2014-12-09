package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("received ctrl+c(%v)\n", sig)
			//os.Exit(0)
		}
	}()

	for {

	}
}
