package main

import (
	"log"
)

func main() {
	var a [4]byte
	//a = make([]byte, 4)
	d := demo(&a)
	log.Printf("%#v\n", a)
	log.Printf("%#v\n", d)
}

func demo(b *[4]byte) [4]byte {
	var c [4]byte
	for i, _ := range b {
		b[i] = 0x67
		c[i] = b[i]
	}

	return c
}
