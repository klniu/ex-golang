package main

import (
    "fmt"
)

func main() {
    a := make([]int, 0)
    fmt.Printf("%#v\n", a)
    a[1] = 2
    fmt.Printf("%#v\n", a)
}