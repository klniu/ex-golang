package main

import (
    "fmt"
    "reflect"
    "time"
)

const times int64  = 10000000

type A struct {
    F int64
}

func (a *A) Add() {
    a.F++
    //fmt.Printf("%#v\n", a.F)
}

func main() {
    test1()
    test2()
    test3()
    test4()
}

func test1() {
    a := &A{}

    begin := time.Now()
    for i := int64(0); i < times; i++ {
        a.Add()
    }
    end := time.Now()
    fmt.Printf("Times(1): %v\n", end.Sub(begin))
}

func test2() {
    a := &A{}

    begin := time.Now()
    for i := int64(0); i < times; i++ {
        val := reflect.ValueOf(a)
        fn := val.MethodByName("Add")
        in := make([]reflect.Value, 0)
        fn.Call(in)
    }
    end := time.Now()
    fmt.Printf("Times(2): %v\n", end.Sub(begin))
}


func test3() {
    a := &A{}

    begin := time.Now()
    val := reflect.ValueOf(a)
    for i := int64(0); i < times; i++ {
        fn := val.MethodByName("Add")
        in := make([]reflect.Value, 0)
        fn.Call(in)
    }
    end := time.Now()
    fmt.Printf("Times(3): %v\n", end.Sub(begin))
}

func test4() {
    a := &A{}

    begin := time.Now()
    val := reflect.ValueOf(a)
    fn := val.MethodByName("Add")
    for i := int64(0); i < times; i++ {
        in := make([]reflect.Value, 0)
        fn.Call(in)
    }
    end := time.Now()
    fmt.Printf("Times(4): %v\n", end.Sub(begin))
}
