package main

import (
	"github.com/zhgo/db"
	"github.com/zhgo/redis/redis"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"time"
)

type Pinyin struct {
	Title     string
	IsPrimary int8
}

type Word struct {
	Title  string
	Url    string
	Pinyin []Pinyin
}

type Zdic struct {
	db       *db.Server
	index    []string
	list     []string
	words    map[string]Word
	retries  int8
	redis    redis.Conn
	getCount int
}

func (z *Zdic) Do() {
	log.SetFlags(log.Lshortfile)

	z.index = make([]string, 0)
	z.list = make([]string, 0)
	z.words = make(map[string]Word)

	db.Env = 3
	z.db = db.Connect("mysql-1", "mysql", "crm_dev:zK8D0krsbrAi@tcp(192.168.22.24:3306)/test?charset=utf8")

	var err error
	z.redis, err = redis.DialTimeout("tcp", "192.168.56.120:6379", 0, 1000*time.Millisecond, 1000*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	begin := time.Now()

	z.Index()

	log.Printf("index: %d\n", len(z.index))

	//sem := make(chan int, 10)

	for _, v := range z.index {
		z.List(v)
	}

	log.Printf("list: %d\n", len(z.list))

	for _, v := range z.list {
		z.Detail(v)
	}

	log.Printf("words: %d\n", len(z.words))
	log.Printf("getCount: %d\n", z.getCount)

	for _, v := range z.words {
		z.Save(v)
	}

	end := time.Now()

	log.Printf("times: %v\n\n", end.Sub(begin))
}

func (z *Zdic) Get(u string) string {
	// Read from redis
	d, err := z.redis.Do("GET", u)
	if err != nil {
		log.Fatal(err)
	}
	if d != nil {
		//log.Fatalf("%#v\n", d)
		return string(d.([]byte))
	}

	z.getCount++

	t := ""
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Referer", "http://www.zdic.net/z/pyjs/")

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("%d %v %v\n", z.retries, resp.Status, u)
		if z.retries < 2 {
			z.retries++
			time.Sleep(1000 * time.Millisecond)
			t = z.Get(u)
		} else {
			log.Fatal(err)
		}
	} else {
		t = string(b)
	}

	// Save to redis
	_, err = z.redis.Do("SET", u, t)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func (z *Zdic) Index() {
	ut := "http://www.zdic.net/z/pyjs/py/?py="
	z.retries = 0
	t := z.Get("http://www.zdic.net/z/pyjs/")

	re := regexp.MustCompile(`shd\(\d+\,[ ]*'(\w+)'\)`)
	matchs := re.FindAllStringSubmatch(t, -1)

	for _, m := range matchs {
		z.index = append(z.index, ut+m[1])
	}
}

func (z *Zdic) List(u string) {
	ut := "http://www.zdic.net"
	z.retries = 0
	t := z.Get(u)

	re := regexp.MustCompile(`HREF="([^"]+)"`)
	matchs := re.FindAllStringSubmatch(t, -1)

	for _, m := range matchs {
		z.list = append(z.list, ut+m[1])
	}
}

func (z *Zdic) Detail(u string) {
	z.retries = 0
	t := z.Get(u)

	title := z.FindStringSubmatch(`“([^”])”字的基本信息`, t, 1)
	pronunciation := z.FindAllStringSubmatch(`<span class="dicpy"><a href="/z/pyjs/\?py=[^"]+" target="_blank">([^<]+)</a><script>spz\("[^"]+"\);<\/script><\/span>`, t, 1)

	word := Word{Title: title, Url: u}

	for i, v := range pronunciation {
		pinyin := Pinyin{Title: v}
		if i == 0 {
			pinyin.IsPrimary = 1
		}
		word.Pinyin = append(word.Pinyin, pinyin)
	}

	z.words[title] = word
}

func (z *Zdic) FindStringSubmatch(exp string, str string, index int) string {
	re := regexp.MustCompile(exp)
	match := re.FindStringSubmatch(str)
	return match[index]
}

func (z *Zdic) FindAllStringSubmatch(exp string, str string, index int) []string {
	re := regexp.MustCompile(exp)
	matchs := re.FindAllStringSubmatch(str, -1)

	ret := make([]string, 0)
	for _, m := range matchs {
		ret = append(ret, m[index])
	}

	return ret
}

func (z *Zdic) Save(w Word) {
	word := db.Item{"title": w.Title, "url": w.Url}
	r, err := z.db.NewQuery().InsertInto("chinese_word").Exec(word)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range w.Pinyin {
		pinyin := db.Item{"word_id": r.LastInsertId, "title": v.Title, "text": v.Title, "is_primary": v.IsPrimary}
		_, err := z.db.NewQuery().InsertInto("chinese_word_pinyin").Exec(pinyin)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	zdic := new(Zdic)
	zdic.Do()
}
