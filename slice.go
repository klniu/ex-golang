package main

import (
	"fmt"
)

func main() {
	var arr [256]byte
	var sli []byte

	arr[0] = 0x41
	sli = arr[0:5]

	fmt.Printf("%#v\n", sli)
	fmt.Printf("%s\n", string(sli))
}