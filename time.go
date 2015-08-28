package main

import (
    "fmt"
    "time"
)

func main() {
    // Mon Jan 2 15:04:05 -0700 MST 2006
    const layout = "2006Mon2 15:04:05"
    t := time.Date(2009, time.November, 10, 15, 0, 0, 0, time.Local)
    fmt.Println(t.Format(layout))
    fmt.Println(t.UTC().Format(layout))
}
