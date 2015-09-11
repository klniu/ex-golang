package main

import(
    "log"
)

func main() {
    a := []byte{0x66, 0x67, 0x68, 0x69, 0x70}
    b := make([]byte, 5)
    c := 0

    for i := 0; i < len(a); i++ {
        b[i] = a[c++]
    }

    log.Printf("%#v\n", b)
}
