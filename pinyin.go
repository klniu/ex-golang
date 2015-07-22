package main

import (
	"bufio"
	"fmt"
	"github.com/zhgo/db"
	"github.com/zhgo/redis/redis"
	"log"
	"os"
	"time"
	"unicode/utf8"
)

var server *db.Server
var cache redis.Conn

func main() {
	log.SetFlags(log.Llongfile)

	db.Env = 3
	server = db.NewServer("mysql-1", "mysql", "crm_dev:zK8D0krsbrAi@tcp(192.168.22.24:3306)/test?charset=utf8")

	var err error
	cache, err = redis.DialTimeout("tcp", "192.168.56.120:6379", 0, 1000*time.Millisecond, 1000*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	UpdateRegister()
	UpdateTag()
}

func UpdateRegister() {
	fw, err := os.Create("pinyin_register.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer fw.Close()

	writer := bufio.NewWriter(fw)

	d := []db.Item{}
	q := server.NewQuery()
	err = q.Select("*").From("ex_service_register").OrderAsc("id").Rows(&d)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range d {
		id := v["id"].(int64)
		province := string(v["province"].([]byte))
		city := string(v["city"].([]byte))
		region := string(v["region"].([]byte))
		log.Printf("%d %d %d %#v %#v %#v\n", id, len(province), utf8.RuneCountInString(province), province, city, region)

		provincePinyin := ""
		for _, j := range province {
			provincePinyin += GetPinyin(string(j), id)
		}
		if province == "陕西" {
			provincePinyin = "shaanxi"
		}
		if province == "重庆" {
			provincePinyin = "chongqing"
		}

		cityPinyin := ""
		for _, j := range city {
			cityPinyin += GetPinyin(string(j), id)
		}
		if city == "重庆" {
			cityPinyin = "chongqing"
		}

		regionPinyin := ""
		for _, j := range region {
			regionPinyin += GetPinyin(string(j), id)
		}

		/*d1 := db.Item{"province_pinyin": provincePinyin, "city_pinyin": cityPinyin, "region_pinyin": regionPinyin}
		q1 := server.NewQuery()
		_, err := q1.Update("ex_service_register").Where(q1.Eq("id", id)).Exec(d1)
		if err != nil {
			log.Fatal(err)
		}*/
		fmt.Fprintf(writer, "UPDATE ex_service_register SET province_pinyin = '%s', city_pinyin = '%s', region_pinyin = '%s' WHERE `id` = %d;\n", provincePinyin, cityPinyin, regionPinyin, id)
	}

	writer.Flush()
}

func UpdateTag() {
	fw, err := os.Create("pinyin_tag.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer fw.Close()

	writer := bufio.NewWriter(fw)

	d := []db.Item{}
	q := server.NewQuery()
	err = q.Select("*").From("ex_tag").Where(q.In("class", 3, 20)).OrderAsc("id").Rows(&d)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range d {
		id := v["id"].(int64)
		class := v["class"].(int64)
		name := string(v["name"].([]byte))
		log.Printf("%d %d %d %#v\n", id, len(name), utf8.RuneCountInString(name), name)

		namePinyin := ""
		for _, j := range name {
			namePinyin += GetPinyin(string(j), id)
		}

		if name == "陕西" {
			namePinyin = "shaanxi"
		}
		if name == "重庆" {
			namePinyin = "chongqing"
		}

		fmt.Fprintf(writer, "UPDATE ex_tag SET url = '%s' WHERE `class` = %d AND `name` = '%s' AND (url = '' OR url IS NULL);\n", namePinyin, class, name)
	}

	writer.Flush()
}

func GetPinyin(char string, id int64) string {
	pinyin := ""

	// Read from redis
	c, err := cache.Do("GET", char)
	if err != nil {
		log.Fatal(err)
	}
	if c != nil {
		pinyin = string(c.([]byte))
		if pinyin == "" {
			log.Printf("%v %v\n", id, char)
		}

		return pinyin
	}

	d := db.Item{}

	q := server.NewQuery()
	err = q.Select("*").From("pinyin").Where(q.Eq("title", char)).Row(&d)
	if err != nil {
		log.Fatal(err)
	}

	if v, ok := d["pinyin"]; ok {
		pinyin = v.(string)
	} else {
		log.Printf("%v %v\n", id, char)
		pinyin = ""
	}

	// Save to redis
	_, err = cache.Do("SET", char, pinyin)
	if err != nil {
		log.Fatal(err)
	}

	return pinyin
}
