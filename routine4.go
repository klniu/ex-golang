package main

import (
    "fmt"
    "time"
)

func main() {
    sem := make(chan int)

    go func() {
        time.Sleep(1 * time.Second)
        fmt.Printf("%s\n", "go in")
        //sem <- 0
    }()

    fmt.Printf("%s\n", "main1")

    <-sem

    fmt.Printf("%s\n", "main2")

    // fatal error: all goroutines are asleep - deadlock!
    // <-sem
}
