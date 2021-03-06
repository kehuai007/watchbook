package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const _zhegx = `<div class="tit"><a href="http://book.zongheng.com/chapter/([\d]+)/([\d]+).html" target="_blank" data-sa-d='{"page_module":"bookDetail","click_name":"newestChapter","book_id":"[\d]+"}'>(.+)</a><em></em></div>`

var regzhBook = regexp.MustCompile(_zhegx)

type ZoneHen struct {
	BookRequest
}

func NewZoneHen(bookRequest BookRequest) *ZoneHen {
	r := &ZoneHen{BookRequest: bookRequest}
	r.BookRequest.Parse = r
	return r
}

func (z ZoneHen) Parse() (book *Book, err error) {
	content, err := get(z.BookUrl)
	if err != nil {
		return nil, err
	}
	match := regzhBook.FindSubmatch(content)
	if len(match) < 4 {
		return nil, errors.New("match failed")
	}
	book = &Book{
		Title: string(match[3]),
		Url:   fmt.Sprintf("http://book.zongheng.com/chapter/%s/%s.html", match[1], match[2]),
		Name:  z.Name,
	}
	book.BookId, _ = strconv.Atoi(string(match[1]))
	book.ChapterId, _ = strconv.Atoi(string(match[2]))
	book.Text = getZhBookText(book.Url)
	fmt.Println(book)
	return
}
func getZhBookText(url string) string {
	text,err := get(url)
	if err != nil {
		return "..."
	}
	match := regexp.MustCompile(`<p>(.+)<p>(.+)</p>`).FindSubmatch(text)
	if len(match) > 0 {
		t := strings.Replace(string(match[0]), "<p>　　", "", -1)
		t = strings.Replace(t, "<p>", "", -1)
		t = strings.Replace(t, "</p>", "", -1)
		return t
	}
	return "..."
}
