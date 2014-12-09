package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		strBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		fmt.Printf("%s\n", string(strBytes))
	}
}
