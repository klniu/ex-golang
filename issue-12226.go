package main

import (
	"fmt"
)

var lines = []string{"foo", "bar", "baz"}

func main() {
	for _, line := range lines {
		if len(line) == 0 ||
			[]byte(line)[0] == []byte("#")[0] {
			continue
		}
		fmt.Println(line)
	}
}
