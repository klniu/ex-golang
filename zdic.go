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
	db        *db.Server
	cache     *redis.Pool
	sliBrowse []string
	sliDetail []string
	words     map[string]Word
	retries   int8
	getCount  int
	routines  int
}

func (z *Zdic) Do() {
	log.SetFlags(log.Lshortfile)

	z.routines = 10
	z.sliBrowse = make([]string, 0)
	z.sliDetail = make([]string, 0)
	z.words = make(map[string]Word)

	db.Env = 3
	z.db = db.Connect("mysql-1", "mysql", "crm_dev:zK8D0krsbrAi@tcp(192.168.22.24:3306)/test?charset=utf8")

	/*var err error
	z.cache, err = redis.DialTimeout("tcp", "192.168.56.120:6379", 0, 1000*time.Millisecond, 1000*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}*/
	z.cache = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "192.168.56.120:6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	begin := time.Now()

	z.Index()
	log.Printf("browse: %d\n", len(z.sliBrowse))

	// Browse
	l := len(z.sliBrowse)
	for i := 0; i < l; i = i + z.routines {
		if i+z.routines >= l {
			z.Browse(z.sliBrowse[i:l])
		} else {
			z.Browse(z.sliBrowse[i : i+z.routines])
		}
	}
	log.Printf("detail: %d\n", len(z.sliDetail))

	// Detail
	l = len(z.sliDetail)
	for i := 0; i < l; i = i + z.routines {
		if i+z.routines >= l {
			z.Detail(z.sliDetail[i:l])
		} else {
			z.Detail(z.sliDetail[i : i+z.routines])
		}
	}
	log.Printf("words: %d\n", len(z.words))
	log.Printf("get: %d\n", z.getCount)

	// Words
	words := make([]Word, 0)
	for _, v := range z.words {
		words = append(words, v)
	}

	l = len(words)
	for i := 0; i < l; i = i + z.routines {
		if i+z.routines >= l {
			z.Save(words[i:l])
		} else {
			z.Save(words[i : i+z.routines])
		}
	}

	end := time.Now()

	log.Printf("times: %v\n\n", end.Sub(begin))
}

func (z *Zdic) Index() {
	z.index("http://www.zdic.net/z/pyjs/", "http://www.zdic.net/z/pyjs/py/?py=")
	z.index("http://www.zdic.net/z/jbs/", "http://www.zdic.net/z/jbs/bs/?bs=")
}

func (z *Zdic) index(u, ut string) {
	z.retries = 0
	t := z.Get(u)

	re := regexp.MustCompile(`shd\(\d+\,[ ]*'([^']+)'\)`)
	matchs := re.FindAllStringSubmatch(t, -1)

	for _, m := range matchs {
		z.sliBrowse = append(z.sliBrowse, ut+m[1])
	}
}

func (z *Zdic) Browse(sli []string) {
	l := len(sli)
	sem := make(chan int, l)

	for i := 0; i < l; i++ {
		go z.browse(sli[i], sem)
	}

	for i := 0; i < l; i++ {
		<-sem
	}
}

func (z *Zdic) browse(u string, sem chan int) {
	ut := "http://www.zdic.net"
	z.retries = 0
	t := z.Get(u)

	re := regexp.MustCompile(`HREF="([^"]+)"`)
	matchs := re.FindAllStringSubmatch(t, -1)

	for _, m := range matchs {
		z.sliDetail = append(z.sliDetail, ut+m[1])
	}

	sem <- 0
}

func (z *Zdic) Detail(sli []string) {
	l := len(sli)
	sem := make(chan int, l)

	for i := 0; i < l; i++ {
		go z.detail(sli[i], sem)
	}

	for i := 0; i < l; i++ {
		<-sem
	}
}

func (z *Zdic) detail(u string, sem chan int) {
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

	sem <- 0
}

func (z *Zdic) Save(sli []Word) {
	l := len(sli)
	sem := make(chan int, l)

	for i := 0; i < l; i++ {
		go z.save(sli[i], sem)
	}

	for i := 0; i < l; i++ {
		<-sem
	}
}

func (z *Zdic) save(w Word, sem chan int) {
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

	sem <- 0
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

func (z *Zdic) Get(u string) string {
	// Pool
	cache := z.cache.Get()
	defer cache.Close()

	// Read from redis
	d, err := cache.Do("GET", u)
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
	_, err = cache.Do("SET", u, t)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	zdic := new(Zdic)
	zdic.Do()
}
