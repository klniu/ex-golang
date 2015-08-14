package main

import (
	"fmt"
)

func main1() {
	a := make([]int, 0)
	fmt.Printf("%#v\n", a)
	a[1] = 2
	fmt.Printf("%#v\n", a)
}

func main() {
	a := make([]string, 0)

	a = append(a, "A")
	a = append(a, "b")
	a = append(a, "C")

	l := len(a)
	for i := 0; i < l; i = i + 2 {
		if i+2 >= l {
			//fmt.Printf("%#v\n", a[i:l-i])
			fmt.Printf("%d %d %#v %#v\n", l, i, a[i:l], a)
		} else {
			fmt.Printf("%d %d %#v %#v\n", l, i, a[i:i+2], a)
		}
	}
}
