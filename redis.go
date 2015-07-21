package main

import (
	"github.com/zhgo/redis/redis"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	r, err := redis.DialTimeout("tcp", "192.168.56.120:6379", 0, 1000*time.Millisecond, 1000*time.Millisecond)

	d, err := r.Do("SET", "AA", "AV")
	if err != nil {
		log.Fatal(err)
	}

	d, err := r.Do("GET", "AA")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", d)

	n, err := redis.Int(c.Do("DBSIZE"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", n)
}
