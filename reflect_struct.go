package main

import (
	"reflect"
	"fmt"
)

type User struct {
	UserId int64 `field: "user_id"`
	Nickname string `field: "nickname"`
}

func main() {
	a := User{}
	b := new(User)

	aType := reflect.TypeOf(a)
	aNumField := aType.NumField()
	for i := 0; i < aNumField; i++ {
		aField := aType.Field(i)
		if !aField.Anonymous && aField.Tag.Get("field") != "" {
			fmt.Printf("%#v\r\n", aField.Tag.Get("field"))
		}
	}

	//方法一
	bType := reflect.ValueOf(b).Type()
	bNumField := bType.Elem().NumField() //注意区别
	for i := 0; i < bNumField; i++ {
		bField := bType.Elem().Field(i) //注意区别
		if !bField.Anonymous && bField.Tag.Get("field") != "" {
			fmt.Printf("%#v\r\n", bField.Tag.Get("field"))
		}
	}

	//方法二
	bType = reflect.ValueOf(b).Elem().Type()
	bNumField = bType.NumField() //注意区别
	for i := 0; i < bNumField; i++ {
		bField := bType.Field(i) //注意区别
		if !bField.Anonymous && bField.Tag.Get("field") != "" {
			fmt.Printf("%#v\r\n", bField.Tag.Get("field"))
		}
	}
}
