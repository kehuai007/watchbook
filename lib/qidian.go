package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const _qdegx = `<a class="blue" href="//(.+)" data-eid="qd_G19" data-cid="//vipreader.qidian.com/chapter/([\d]+)/([\d]+)" title="(.+)" target="_blank">.+</a><i>`
const _qdegxtext =`<p>(.+)<p>(.+)`
var regqdBook = regexp.MustCompile(_qdegx)
var regqdText = regexp.MustCompile(_qdegxtext)
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
		Title: strings.TrimSpace(string(match[4])),
		Url:   "https://" + string(match[1]),
		Name:  strings.TrimSpace(q.Name),
	}
	book.BookId, _ = strconv.Atoi(string(match[2]))
	book.ChapterId, _ = strconv.Atoi(string(match[3]))
	book.Text = getQdBookText(book.Url)
	fmt.Println(book)
	return
}
func getQdBookText(url string) string  {
	text,err := get(url)
	if err != nil {
		return "..."
	}
	match := regqdText.FindSubmatch(text)
	if len(match) == 0{
		match = regexp.MustCompile(`<p>(.+)<p>`).FindSubmatch(text)
	}
	if len(match) > 0 {
		t := strings.Replace(string(match[0]), "<p>　　", "", -1)
		t = strings.Replace(t, "<p>", "", -1)
		return t+"..."
	}
	return "..."
}
