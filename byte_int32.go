package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

func main() {
	b := []byte{0xd0, 0xb4}
	b_buf := bytes.NewBuffer(b)
	var x int32
	binary.Read(b_buf, binary.BigEndian, &x)
	fmt.Println(x)

	fmt.Println(strings.Repeat("-", 100))

	x = 1076
	b_buf = bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, x)
	fmt.Println(b_buf.Bytes())
}
