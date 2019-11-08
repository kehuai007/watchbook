package main

import (
	"log"
	"mall/wxpush"
	"mall/zhwatch"
)

func SendSubMessage(lastBook *zhwatch.BookInfo, alreadySend, avUsers map[string]bool, bookChan <-chan *zhwatch.BookInfo, client RedisConn) {
	for {
		select {
		case book := <-bookChan:
			if lastBook == nil || book.Index != lastBook.Index {
				lastBook = book
				client.SetLastBook(lastBookKey, book)
				alreadySend = make(map[string]bool)
				client.Del(alreadyKey)
			}
			readySend := make([]string, 0)
			for usr, enable := range avUsers {
				if _, ok := alreadySend[usr]; !ok && enable {
					readySend = append(readySend, usr)
				}
				SendBookInfo(book, readySend, alreadySend, client)
			}
		}
	}
}

//> 发送书本最新章节
func SendBookInfo(book *zhwatch.BookInfo, uid []string, alreadySend map[string]bool, client RedisConn) error {
	msg := wxpush.NewMessage(appToken)
	msg.UIds = append(msg.UIds, uid...)
	if len(msg.UIds) == 0 {
		return wxpush.NewError(1008, nil)
	}
	msg.SetContent("《逆天邪神》更新啦，" + book.GetBookInfo())
	msg.SetUrl(book.Url)
	result, err := wxpush.SendMessage(*msg)
	if err != nil {
		log.Println("SendMgs err", err)
		return err
	}
	for _, u := range msg.UIds {
		alreadySend[u] = true
		client.SetAlready(alreadyKey, u)
	}
	log.Println(result)
	return nil
}
