package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	//fileInput, err := os.Open("input.txt")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//fileOutput, err := os.Open("output.txt")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//fileError, err := os.Open("error.txt")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	var procAttr os.ProcAttr

	//procAttr.Files = []*os.File{fileInput, fileOutput, fileError}
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	//procAttr.Sys.HideWindow = true

	_, err := os.StartProcess("ls.exe", nil, &procAttr)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(p.Pid)
}

func cmd(str string, arg ...string) {
	out, err := exec.Command(str, arg...).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", out)
}
