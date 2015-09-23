package main

import (
    "fmt"
)

func main() {
    sem := make(chan int)

    // 1
    // fatal error: all goroutines are asleep - deadlock!
    //sem <- 0

    // 2
    // fatal error: all goroutines are asleep - deadlock!
    //<-sem

    // 3
    // fatal error: all goroutines are asleep - deadlock!
    sem <- 0
    s := <-sem
    fmt.Printf("%d\n", s)
}
