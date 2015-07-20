package main

import(
	"github.com/zhgo/db"
	"log"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Pinyin struct {
	Title string 
	IsPrimary int8
}

type Word struct {
	Title string
	Pinyin []Pinyin
}

type Zdic struct {
	server  *db.Server
	index   []string
	list    []string
	words   map[string]Word
}

func (z *Zdic) Do() {
	log.SetFlags(log.Lshortfile)

	z.index = make([]string, 0)
	z.list  = make([]string, 0)
	z.words  = make(map[string]Word)

	z.server = db.NewServer("mysql-1", "mysql", "crm_dev:zK8D0krsbrAi@tcp(192.168.22.24:3306)/test?charset=utf8")

	begin := time.Now()

	z.Index()
	
	for _, v := range z.index {
		z.List(v)
	}

	for _, v := range z.list {
		z.Detail(v)
	}
	
	for _, v := range z.words {
		z.Save(v)
	}

	end := time.Now()
	
	log.Printf("index: %d\n", len(z.index))
	log.Printf("list: %d\n", len(z.list))
	log.Printf("words: %d\n", len(z.words))
	log.Printf("times: %v\n", end.Sub(begin))
}

func (z *Zdic) Get(u string) string {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Referer", "http://www.zdic.net/z/pyjs/")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}

func (z *Zdic) Index() {
	ut := "http://www.zdic.net/z/pyjs/py/?py="
	t := z.Get("http://www.zdic.net/z/pyjs/")

	re := regexp.MustCompile(`shd\(\d+\,[ ]*'(\w+)'\)`)
	matchs := re.FindAllStringSubmatch(t, -1)

	for _, m := range matchs {
		z.index = append(z.index, ut + m[1])
	}
}

func (z *Zdic) List(u string) {
	ut := "http://www.zdic.net"
	t := z.Get(u)

	re := regexp.MustCompile(`HREF="([^"]+)"`)
	matchs := re.FindAllStringSubmatch(t, -1)

	for _, m := range matchs {
		z.list = append(z.list, ut + m[1])
	}
}

func (z *Zdic) Detail(u string) {
	t := z.Get(u)

	title := z.FindStringSubmatch(`“([^”])”字的基本信息`, t, 1)
	pronunciation := z.FindAllStringSubmatch(`<span class="dicpy"><a href="[^"]+" target="_blank">([^<]+)</a><script>spz\("[^"]+"\);<\/script><\/span>`, t, 1)

	word := Word{Title: title}

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
	word := db.Item{"title": w.Title}
	r, err := db.NewQuery(z.server).InsertInto("chinese_word").Exec(word)
    if err != nil {
        log.Fatal(err)
    }

    for _, v := range w.Pinyin {
    	pinyin := db.Item{"word_id": r.LastInsertId, "title": v.Title, "is_primary": v.IsPrimary}
		_, err := db.NewQuery(z.server).InsertInto("chinese_word_pinyin").Exec(pinyin)
	    if err != nil {
	        log.Fatal(err)
	    }
    }
}

func main() {
	zdic := new(Zdic)
	zdic.Do()
}