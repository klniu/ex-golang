package main

import (
	"fmt"
)

const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Partyday
	numberOfDays  // this constant is not exported
)

func main() {
	fmt.Printf("%T", Sunday)
	fmt.Println(Monday)
}
