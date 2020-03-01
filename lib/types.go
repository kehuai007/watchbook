package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type BookRequest struct {
	BookUrl string
	Name    string
	BookId  int
	Parse   Parser
	Closed  chan bool
}

func (b *BookRequest) GetBookRequest() BookRequest {
	return *b
}
func (b *BookRequest) Run(minute int) chan *Book {
	ch := make(chan *Book)
	book, _ := b.Parse.Parse()
	tick := time.Tick(time.Minute * time.Duration(minute))
	go func() {
		for {
			select {
			case <-tick:
				b, err := b.Parse.Parse()
				if err != nil {
					continue
				}
				if book.ChapterId != b.ChapterId {
					book = b
					ch <- b
				}
			case <-b.Closed:
				log.Println("exit by chan")
				return
			}
		}
	}()
	return ch
}

// 书籍信息
type Book struct {
	ChapterId int    `json:"chapterid"`
	BookId    int    `json:"bookid"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	Name      string `json:"name"`
	Text      string `json:"text"`
}
type Parser interface {
	Parse() (*Book, error)
}
type WatchBookServer interface {
	Run(int) chan *Book
	GetBookRequest() BookRequest
}

// get method
func get(url string) (buf []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
