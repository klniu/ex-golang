package main

import (
	"log"
	"fmt"
	"errors"
)

func main() {
	//a, err都未定义, 允许使用:=
	a, err := funSuccess()
	if err != nil {
		log.Fatal(err)
	}

	//a已定义, err1未定义 ?
	a, err1 := funSuccess()
	if err1 != nil {
		log.Fatal(err1)
	}

	//err已经定义, 但b未定义, 允许使用:=
	b, err := funFailure()
	if err != nil {
		log.Fatal(err)
	}

	//err := fun1() //err已经定义, 且不存在其他未定义变量, 不能使用:=
	err = fun1()
	if err != nil {
		log.Fatal(err)
	}

	//err := fun2() //同上
	err = fun2()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(a)
	fmt.Println(b)
}

func funSuccess() (string, error) {
	return "Success", nil
}

func funFailure() (string, error) {
	return "", errors.New("Failure")
}

func fun1() error {
	return nil
}

func fun2() error {
	return errors.New("fun2 failure")
}
