package main

import (
    "fmt"
    "regexp"
    "strconv"
    "log"
)

func main1() {
    re := regexp.MustCompile("a(x*)b")
    fmt.Println(re.ReplaceAllString("-ab-axxb-", "T"))
    fmt.Println(re.ReplaceAllString("-ab-axxb-", "$1"))
    fmt.Println(re.ReplaceAllString("-ab-axxb-", "$1W"))
    fmt.Println(re.ReplaceAllString("-ab-axxb-", "${1}W"))
}

func main2() {
    re := regexp.MustCompile("(a.)b")
    fmt.Println(re.FindAllString("parbanormalb", -1))
    fmt.Println(re.FindAllString("paranbormal", 2))
    fmt.Println(re.FindAllString("graalb", -1))
    fmt.Println(re.FindAllString("none", -1))
}

func main() {
    re := regexp.MustCompile("(a.)b")
    fmt.Println(re.FindAllStringSubmatch("parbanormalb", -1))
    fmt.Println(re.FindAllStringSubmatch("paranbormal", 2))
    fmt.Println(re.FindAllStringSubmatch("graalb", -1))
    fmt.Println(re.FindAllStringSubmatch("none", -1))
}

func main4() {
    str := "INSERT INTO table1 (c1, c2, c3, c4, c5) VALUES ($1, $2, $3, $2, $3)"
    args := []interface{}{1980, "Bob", "Male"}

    re := regexp.MustCompile(`\$(\d+)`)
    sli := re.FindAllStringSubmatch(str, -1)
    //fmt.Printf("%#v\n\n", sli)

    newArgs := make([]interface{}, len(sli))
    for i, v := range sli {
        vi, err := strconv.ParseInt(v[1], 10, 0)
        if err != nil {
            log.Fatalf("%v\n", err)
        }
        newArgs[i] = args[vi-1]
    }

    newStr := re.ReplaceAllString(str, "?")

    fmt.Printf("%#v\n\n", newStr)
    fmt.Printf("%#v\n\n", newArgs)

}

func main5() {
    str := "INSERT INTO table1 (f1) VALUES('Bob\\'s Book')"


    re := regexp.MustCompile(`'[^']+'`)
    fmt.Printf("%#v\n", re.FindAllString(str, -1))
}