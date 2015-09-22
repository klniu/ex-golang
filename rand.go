// go run rand.go
// Author: liudng@gmail.com
// 2015-9-22

package main

import (
    "fmt"
    "log"
    "math/rand"
    "os"
    "runtime/pprof"
    "strings"
    "time"
)

func main() {
    fc, err := os.Create("rand.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer fc.Close()
    pprof.StartCPUProfile(fc)
    defer pprof.StopCPUProfile()

    fm, err := os.Create("rand.mprof")
    if err != nil {
        log.Fatal(err)
    }
    pprof.WriteHeapProfile(fm)
    fm.Close()

    //Begin
    begin := time.Now()

    totalStr := 600 * 10000
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

    end := time.Now()

    fmt.Printf("共计%v万条数据，数据总长度%v\n", totalStr/10000, totalLength)
    fmt.Printf("完成过滤共耗时%v\n", end.Sub(begin))
}
