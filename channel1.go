package main

import (
    "fmt"
    "time"
)

func sleepAndTalk(secs time.Duration, msg string, c chan string) {
    time.Sleep(secs * time.Second)
    c <- msg
}

func main() {
    c := make(chan string)

    go sleepAndTalk(0, "Hello", c)
    go sleepAndTalk(1, "Gophers!", c)
    go sleepAndTalk(2, "What's", c)
    go sleepAndTalk(3, "up?", c)

    for i := 0; i < 4; i++ {
        printChan(c)
        fmt.Printf("%#v\n", <-c)
        printChan(c)
    }
}

func printChan(v chan string) {
    for i := 0; i < len(v); i++ {
        fmt.Printf("%#v\n", v[i])
    }
}
