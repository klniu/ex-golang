package main

import (
    "fmt"
    "net"
)

const (
    addr = "127.0.0.1:80"
)

func main() {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        fmt.Println("Connet failed:", err.Error())
        return
    }
    fmt.Println("Connet Success")
    defer conn.Close()

    Client(conn)
}

func Client(conn net.Conn) {
    sms := make([]byte, 1024)
    for {
        fmt.Print("Please input:")
        _, err := fmt.Scan(&sms)
        if err != nil {
            fmt.Println("Input error:", err.Error())
        }

        h := []byte("GET http://127.0.0.1/ HTTP/1.1\nHost: 127.0.0.1\n")
        for _, v := range sms {
            h = append(h, v)
        }
        h = append(h, byte('\n'))

        conn.Write(h)

        buf := make([]byte, 1024)
        c, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Response failed:", err.Error())
        }
        fmt.Println(string(buf[0:c]))
    }

}
