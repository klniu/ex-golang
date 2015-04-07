package main

import (
    "fmt"
)

func main() {
    a := map[string]string{"A": "a", "B": "b"}
    b := map[string]string{"B": "Bb", "C": "c"}

    c := append(a, b)

    fmt.Printf("%#v", c)
}
