package main

import (
    "fmt"
    "time"
)

func sleepAndTalk(secs time.Duration, msg string, c chan string) {
    c <- msg
    time.Sleep(secs * time.Second)
}

func main() {
    c := make(chan string, 6)

    go sleepAndTalk(0, "Hello", c)
    go sleepAndTalk(1, "Gophers!", c)
    go sleepAndTalk(2, "What's", c)
    go sleepAndTalk(3, "up?", c)

    for i := 0; i < 4; i++ {
        fmt.Printf("%v ", <-c)
    }
}
