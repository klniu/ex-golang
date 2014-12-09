package main

import(
	"bytes"
	//"encoding/binary"
	"fmt"
)

func main() {
	num()
}

func str() {
	a := "我"
	b := []byte(a)
	//c := []int8(a) // Wrong convert
	//c1 := []uint8(a)
	//d := []int16(a) // Wrong convert
	//d1 := []uint16(a) // Wrong convert
	e := []rune(a)
	//e1 := []uint32(a) // Wrong convert
	//f := []int64(a) // Wrong convert
	//f1 := []uint64(a) // Wrong convert

	//fmt.Println(a, b, c, c1, d, d1, e, e1, f, f1)
	//fmt.Println(a, b, c1, e)
	fmt.Printf("%#v\n", a)
	fmt.Printf("%+q\n", a)
	fmt.Printf("%#v\n", b)
	//fmt.Printf("%#v\n", c1)
	fmt.Printf("%#v\n", e)
	fmt.Printf("%#v\n", bytes.Runes(b))

	//for i := 0; i < len(a); i++ {
    //    fmt.Printf("%#v ", a[i])
    //}

    //b_buf := bytes.NewBuffer(b)
	//var x int32
	//binary.Read(b_buf, binary.BigEndian, &x)
	//fmt.Println(x)

	//b_buf = bytes.NewBuffer([]byte{})
	//binary.Write(b_buf, binary.BigEndian, e)
	//fmt.Println(b_buf.Bytes())
}

func num() {
	a := 8
	b := []byte(a)
	e := []rune(a)
	
	fmt.Printf("%#v\n", a)
	fmt.Printf("%+q\n", a)
	fmt.Printf("%#v\n", b)
	fmt.Printf("%#v\n", e)
	fmt.Printf("%#v\n", bytes.Runes(b))
}