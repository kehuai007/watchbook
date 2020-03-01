package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const _qdegx = `<a class="blue" href="//(.+)" data-eid="qd_G19" data-cid="//vipreader.qidian.com/chapter/([\d]+)/([\d]+)" title="(.+)" target="_blank">.+</a><i>`

var regqdBook = regexp.MustCompile(_qdegx)

type QiDian struct {
	BookRequest
}

func NewQiDian(bookRequest BookRequest) *QiDian {
	r := &QiDian{BookRequest: bookRequest}
	r.BookRequest.Parse = r
	return r
}

func (q QiDian) Parse() (book *Book, err error) {
	content, err := get(q.BookUrl)
	if err != nil {
		return nil, err
	}
	match := regqdBook.FindSubmatch(content)
	if len(match) < 5 {
		return nil, errors.New("match failed")
	}
	book = &Book{
		Title: string(match[4]),
		Url:   "https://" + string(match[1]),
		Name:  q.Name,
	}
	book.BookId, _ = strconv.Atoi(string(match[2]))
	book.ChapterId, _ = strconv.Atoi(string(match[3]))
	fmt.Println(book)
	return
}
