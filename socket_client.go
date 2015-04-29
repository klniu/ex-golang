package main

import (
    "bufio"
    "fmt"
    "net"
    "log"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:80")
    if err != nil {
        log.Fatalf("Connection error:%s", err)
    }

    connbuf := bufio.NewReader(conn)
    for{
        str, err := connbuf.ReadString('\n')
        if err!= nil {
            log.Fatal(err)
        }
        fmt.Println("Respones: ")
        if len(str)>0 {
            fmt.Println(str)
        }
    }
}