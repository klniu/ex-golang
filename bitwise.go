package main

import (
    "log"
)

func main() {
    var a, b, c, d uint32
    a = 0x67
    b = 0x56
    c = 0x88
    d = 0x27

    x := uint32(a | b)
    y := a | b | c

    log.Printf("%#x\n", x)
    log.Printf("%#x\n", y)

    m := x & a
    n := x & d
    log.Printf("%#x\n", m)
    log.Printf("%#x\n", n)

    if (x & a) != 0 {
        log.Printf("%#x\n", x&a)
    }
    if (x & d) != 0 {
        log.Printf("%#x\n", x&d)
    }
}
