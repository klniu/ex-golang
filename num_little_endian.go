package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	array1 := []byte{0xff, 0x00, 0x00}
	num1, _ := binary.Uvarint(array1)
	fmt.Printf("%v, %x, %v\n", array1, num1, num1)

	num2 := unpackNumber(array1, 8)
	fmt.Printf("%v, %x, %v\n", array1, num2, num2)
}

//Convert n bytes to uint64 (Little Endian)
func unpackNumber(b []byte, n uint8) uint64 {
	if n < 1 {
		return 0
	}
	var r uint64 = 0
	for i := uint8(0); i < n; i++ {
		r |= (uint64(b[i]) << (i * 8))
	}
	return r
}
