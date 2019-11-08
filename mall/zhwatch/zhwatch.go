package zhwatch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const QueryUrl = "http://book.zongheng.com/book/${bookId}.html"

type BookInfo struct {
	Title     string `json:"title"`
	Url       string `json:"url"`
	BookId    string `json:"bookid"`
	ArticleId string `json:"articleid"`
	Index     string `json:"index"`
}

func (b *BookInfo) GetBookInfo() string {
	s := "最新第" + b.Index + "章 " + b.Title + " "
	return s
}
func QueryBookInfo(bookId int) (*BookInfo, error) {
	url := strings.ReplaceAll(QueryUrl, "${bookId}", strconv.Itoa(bookId))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	reg := regexp.MustCompile(`<div class="tit"><a .*?href="(.*?)".*?>.*?</a>`)
	matched := reg.Find(respBody)
	var result = BookInfo{}
	s := string(matched)
	first := strings.Index(s, `http://`)
	last := strings.Index(s, `.html`)
	result.Url = s[first : last+len(`.html`)]
	first = strings.LastIndex(result.Url, `/`)
	last = strings.Index(result.Url, `.html`)
	result.ArticleId = result.Url[first+len(`/`) : last]
	first = strings.Index(s, `"book_id":"`)
	last = strings.Index(s, `"}'>`)
	result.BookId = s[first+len(`"book_id":"`) : last]
	first = last
	last = strings.Index(s, `</a>`)
	title := s[first+len(`"}'>`) : last]
	first = strings.Index(title, `第`)
	last = strings.Index(title, `章`)
	result.Index = title[first+len(`第`) : last]
	result.Title = strings.TrimSpace(title[last+len(`章`):])
	return &result, nil
}
func GetBook(bookId int) ([]byte, error) {
	book, err := QueryBookInfo(408586)
	if err != nil {
		return nil, err
	}
	s, err := json.Marshal(book)
	if err != nil {
		return nil, err
	}
	return s, nil
}
