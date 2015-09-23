package main

import (
    "fmt"
    "time"
)

func main() {
    sem := make(chan int, 20)

    // spend ~2 seconds.
    for i := 0; i < 50; i++ {
        fmt.Printf("i %d\n", i)
        go func(s int) {
            sem <- s
            time.Sleep(1 * time.Second)
            fmt.Printf("in %d\n", s)
            <-sem
        }(i)
    }

    // wait 10 seconds, haha.
    time.Sleep(10 * time.Second)
}
