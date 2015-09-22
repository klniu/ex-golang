// go run hash.go
// Author: liudng@gmail.com
// 2014-9-23

package main

import (
    "fmt"
    "math/rand"
    "sort"
    "strings"
    "time"
)

type sortBytes []byte

func (s sortBytes) Less(i, j int) bool {
    return s[i] < s[j]
}

func (s sortBytes) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s sortBytes) Len() int {
    return len(s)
}

func main() {
    totalStr := 200 * 10000
    myArr := make([]string, totalStr)
    totalLength := 0
    input := strings.Split("白日依山尽黄河入海流欲穷千里目更上一层楼危楼高百尺可以摘星辰不感高声语恐惊天上人", "")

    //生成随机字符串, 长度为3~16
    rand.Seed(time.Now().Unix())
    for i := 0; i < totalStr; i++ {
        tempStr := ""
        tempLength := rand.Intn(13) + 3
        for j := 0; j < tempLength; j++ {
            tempStr += input[rand.Intn(len(input))]
        }
        totalLength += tempLength
        myArr[i] = tempStr
    }

    //Begin
    begin := time.Now()

    dict := make(map[string]string)
    for _, v := range myArr {
        keySli := []byte(v)
        sort.Sort(sortBytes(keySli))
        key := string(keySli)
        if _, s := dict[key]; s == false {
            dict[key] = v
        }
    }

    end := time.Now()

    fmt.Printf("共计%v万条数据，数据总长度%v, 其中%v条不重复数据\n", totalStr/10000, totalLength, len(dict))
    //fmt.Printf("%v\n%v\n", begin, end)
    fmt.Printf("完成过滤共耗时%v\n", end.Sub(begin))
}
