// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//output to file
//var fileOutput *os.File
//var fileOutputErr error

//config stract
type Config struct {
	MplusData string                 `json:"mplus_data"`
	Programs  map[string]confProgram `json:"programs"`
	Chdir     map[string]string      `json:"chdir"`
	Env       map[string]string      `json:"env"`
	Path      []string               `json:"path"`
	AutoRuns  map[string]confService `json:"auto_runs"`
	Services  map[string]confService `json:"services"`
}

type confProgram struct {
	Typ  int8          `json:"type"`
	Path string        `json:"path"`
	Args []interface{} `json:"args"`
}

type confService struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

//Config
var config Config

//root path
var basePath string

func main() {
	//init output to file
	//fileOutput, fileOutputErr = os.Create("output.txt")
	//if fileOutputErr != nil {
	//	log.Fatalln(fileOutputErr)
	//}

	//init basePath
	var err error
	//basePath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	basePath, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//Load config file
	initConfig()

	//fmt.Printf("%#v\n", config)

	//autoRun()

	//kill all running services
	serviceStop("all")

	//start all service at this time
	serviceStart("all")

	//file monitor
	//go fileMonitor()

	//console cycle
	console()

	//kill all service at exit
	serviceStop("all")
}

func initConfig() {
	txt, err := ioutil.ReadFile(basePath + "/mplus.json")
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(strings.NewReader(string(txt)))
	err = dec.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	//Dir alias
	for k, v := range config.Chdir {
		config.Chdir[k] = repx(v)
	}

	//Env
	setenv("MPLUS_DATA", config.MplusData)
	for k, v := range config.Env {
		config.Env[k] = repx(v)
		setenv(k, config.Env[k])
	}

	//Path
	path := ""
	for i, v := range config.Path {
		config.Path[i] = repx(v)
		path += config.Path[i] + ";"
	}
	setenv("PATH", path)

	//Services
	for k, v := range config.Services {
		args := v.Args

		for i, s := range args {
			a := repx(s)
			args[i] = a
		}

		v.Args = args
		config.Services[k] = v
	}
}

func repx(str string) string {
	re := regexp.MustCompile(`\{basePath\}`)
	str = re.ReplaceAllLiteralString(str, basePath)

	re = regexp.MustCompile(`\{mplusData\}`)
	str = re.ReplaceAllLiteralString(str, config.MplusData)

	re = regexp.MustCompile(`\{SystemRoot\}`)
	str = re.ReplaceAllLiteralString(str, os.Getenv("SystemRoot"))

	re = regexp.MustCompile(`\{UserProfile\}`)
	str = re.ReplaceAllLiteralString(str, os.Getenv("USERPROFILE"))

	return str
}

func console() {
ConsoleLoop:
	for {
		path, errPath := os.Getwd()
		if errPath != nil {
			fmt.Printf("%s\n", errPath)
		}

		fmt.Print("[" + path + "] ")

		reader := bufio.NewReader(os.Stdin)
		strBytes, _, errReadLine := reader.ReadLine()
		if errReadLine != nil {
			fmt.Printf("%s\n", errReadLine)
		}

		args := strings.Split(string(strBytes), " ")

		switch args[0] {
		case "cd": //change dir command
			chdir(args[1])
		case "dir", "cls", "del", "deltree", "path", "set": //inner command
			arg := append([]string{"/C"}, args...)
			exe(2, "cmd.exe", arg...)
		case "q", "quit", "exit": //exit mplus
			break ConsoleLoop
		case "~", "home": //mplus home dir
			chdir(basePath)
		case "service": //start or stop: nginx, php-cgi, mysql, memcached
			if len(args) > 2 {
				if args[1] == "stop" {
					if args[2] == "all" {
						serviceStop("all")
					} else {
						serviceStop(args[2])
					}
				} else if args[1] == "start" {
					if args[2] == "all" {
						serviceStart("all")
					} else {
						serviceStart(args[2])
					}
				}
			}
		default:
			//find Programs
			cmd1, s1 := config.Programs[args[0]]
			if s1 == true {
				exeAny(cmd1.Typ, cmd1.Path, cmd1.Args...)
				continue
			}

			//find Chdir
			cmd2, s2 := config.Chdir[args[0]]
			if s2 == true {
				chdir(cmd2)
				continue
			}

			exe(2, args[0], args[1:]...)
		}
	}
}

func exeAny(typ int8, name string, args ...interface{}) *exec.Cmd {
	a := make([]string, len(args))

	for i, v := range args {
		a[i] = v.(string)
	}

	return exe(typ, name, a...)
}

func exe(typ int8, name string, args ...string) *exec.Cmd {
	if name == "new" && len(args) > 0 {
		name = "cmd.exe"
		args = append([]string{"/C", "start"}, args...)
	}

	cmd := exec.Command(name, args...)

	switch typ {
	case 1: //service mode
		fmt.Printf("[%s]", name)
		fmt.Println(args)

		err := cmd.Start()
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		fmt.Println(cmd.Process.Pid)

		//runtime.SetFinalizer(cmd, cmd.Process.Kill())
	case 2: //application mode, input or output to standard window
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}

	return cmd
}

func setenv(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func chdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func autoRun() {
	//  git config --global user.email "liudng@gmail.com"
	//  git config --global user.name "liudng"
	for _, v := range config.AutoRuns {
		exe(2, v.Cmd, v.Args...)
	}
}

func serviceStart(name string) {
	if name == "all" {
		for k, _ := range config.Services {
			serviceStart(k)
		}

		return
	}

	service, s := config.Services[name]
	if s == false {
		fmt.Println("Service not existing")
		return
	}

	if service.Cmd == "nginx.exe" {
		//delete error log
		os.Remove(fmt.Sprintf("%s/nginx_x86/logs/error.log", basePath))
		chdir(basePath + "/nginx_x86")
	}

	exe(1, service.Cmd, service.Args...)

	if service.Cmd == "nginx.exe" {
		chdir(basePath)
	}
}

func serviceStop(name string) {
	if name == "all" {
		for k, _ := range config.Services {
			serviceStop(k)
		}

		return
	}

	service, s := config.Services[name]
	if s == false {
		fmt.Println("Service not existing")
		return
	}

	/*if name == "redis" {
		exe(1, "redis-cli.exe", "shutdown")
	}*/

	exe(2, "taskkill", "/f", "/im", service.Cmd)
}
