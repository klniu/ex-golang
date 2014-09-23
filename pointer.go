package main

import (
	"reflect"
	"fmt"
)

func main() {
	data := []int{1, 2}
	dataPtr := &data

	aaa := []int{7, 8, 9}

	//直接赋值: 只要是同类型的切片, 可互相赋值
	fmt.Printf("%#v\r\n", data)
	data = aaa
	fmt.Printf("%#v\r\n\r\n", data)

	//改变指针dataPtr的值, data的值并没有改变
	fmt.Printf("%#v\r\n", data)
	fmt.Printf("%#v\r\n", dataPtr)
	dataPtr = &aaa
	fmt.Printf("%#v\r\n", data) //改变dataPtr的值, 不会改变data
	fmt.Printf("%#v\r\n\r\n", dataPtr) //dataPtr指向的值已经发生改变

	//fmt.Printf("%#v\r\n\r\n", data)
	//fun1(dataPtr)
	//fmt.Printf("%#v\r\n\r\n", data)

	//fmt.Printf("%#v\r\n\r\n", data)
	//fun2(dataPtr)
	//fmt.Printf("%#v\r\n\r\n", data)

	//fmt.Printf("%#v\r\n\r\n", data)
	//fun3(dataPtr)
	//fmt.Printf("%#v\r\n\r\n", data)

}

//例1: 通指针为传入指针赋值
func fun1(ptr *[]int) {
	aaa := []int{4, 5, 6}
	fmt.Printf("%#v\r\n\r\n", ptr)
	*ptr = aaa //改变
	fmt.Printf("%#v\r\n\r\n", ptr)

}

//例2: 通指针为传入指针赋值, 无法实现
//func fun3(ptr *[]interface{}) {
func fun2(ptr interface{}) {
	aaa := []int{4, 5, 6}
	fmt.Printf("%#v\r\n\r\n", ptr)
	ptr = &aaa //ptr和dataPtr是不同的变量, 都指向data的指针, 改动ptr的指向不会改变dataPtr, 也不会改变data的值
	fmt.Printf("%#v\r\n\r\n", ptr)

	//main.go:27: cannot use dataPtr (type *[]int) as type *[]interface {} in argument to fun3
	//main.go:57: cannot use aaa (type []int) as type []interface {} in assignment exit status 2
}

//例3: 通过refrect给传入数组参数赋值
func fun3(ptr interface{}) {
	aaa := 6 //自定义值
	valuePtr := reflect.Indirect(reflect.ValueOf(ptr)) //ptr的Value
	valuePtr.Set(reflect.Append(valuePtr, reflect.Indirect(reflect.ValueOf(&aaa)))) //赋值
}
