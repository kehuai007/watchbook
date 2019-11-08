package main

import (
	"log"
	"mall/wxpush"
	"mall/zhwatch"
)

const (
	//> 发消息获取用户必须要的appToken
	appToken = `AT_****`
	//> 纵横网书籍 逆天邪神
	bookId = 408586
	//> 用户注册回调监听端口
	recvUserPort = 8080
	//> 获取图书信息时间间隔 单位分钟
	timerMinute = 2
	//> alreadyKey
	prefix      = `AT_mISo`
	alreadyKey  = prefix + `already`
	lastBookKey = prefix + `lastBook`
	redisIp     = `127.0.01`
	redisPort   = `6379`
)

func main() {
	//> 所有订阅用户
	allUsers := make(map[string]bool, 0)
	//> 订阅已经发送的用户
	alreadySend := make(map[string]bool)
	client := NewRedisConn(redisIp, redisPort)
	client.GetAllAlready(alreadyKey, alreadySend)
	//> 书籍队列
	bookChan := make(chan *zhwatch.BookInfo)
	//> 最新一章的书籍
	var lastBook *zhwatch.BookInfo
	client.GetLastBook(lastBookKey, lastBook)
	//> 新增订阅用户
	recvUsrChan := make(chan *wxpush.RecvUserData)

	//> 获取全部用户
	GetAllUser(allUsers)
	//> 获取书籍
	go QueryBookInfo(bookChan)
	//> 发送书籍信息
	go SendSubMessage(lastBook, alreadySend, allUsers, bookChan, client)
	//> 监听关注用户
	go OnUserSub(lastBook, alreadySend, allUsers, recvUsrChan, client)
	log.Println("Server Start....")
	//> 开启监听
	wxpush.ServerAndListen(recvUserPort, recvUsrChan)
}
