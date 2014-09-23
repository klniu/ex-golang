package main

import (
	"fmt"
	"log"
	"os"
)

//批量生成逗号分隔的数字字符串
func main() {
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}

	// close fo on exit and check for its returned error
	defer f.Close()

	var l int64
	l = 0
	for i := 50001; i <= 50030; i++ {
		b := []byte(fmt.Sprintf("%d,", i))

		n, err := f.WriteAt(b, l)
		if err != nil {
			log.Fatal(err)
		}

		l = l + int64(n)

		fmt.Println(n, l, i)
	}
}
