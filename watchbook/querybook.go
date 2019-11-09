package main

import (
	"log"
	"mall/zhwatch"
	"time"
)

//> 定时获取书籍
func QueryBookInfo(bookChan chan<- *zhwatch.BookInfo) {
	t := time.Tick(time.Minute * timerMinute)
	for {
		select {
		case <-t:
			book, err := zhwatch.QueryBookInfo(bookId)
			if err != nil {
				log.Println("QueryBookInfo err ", err)
				continue
			}
			bookChan <- book
		}
	}
}
