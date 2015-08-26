package main

import (
	"encoding/binary"
	"log"
)

var a []byte = []byte{0xff, 0xf7}
var b []byte = []byte{0x0f, 0x80}

func main() {
	// ff f7
	// 0f 80

	method1()
	method2()
}

func method1() {
	x := uint32(binary.LittleEndian.Uint16(a))
	y := uint32(binary.LittleEndian.Uint16(b))<<16 | x
	log.Printf("%#v\n", x)
	log.Printf("%#v\n", y)
}

func method2() {
	x := uint32(a[0]) | uint32(a[1])<<8
	y := uint32(a[0]) | uint32(a[1])<<8 | uint32(b[1])<<24 | uint32(b[0])<<16
	log.Printf("%#v\n", x)
	log.Printf("%#v\n", y)
}
