package main

import (
	"log"
	"mall/wxpush"
	"mall/zhwatch"
)

func GetAllUser(avUsers map[string]bool) error {
	result, err := wxpush.QueryAllUser(appToken)
	if err != nil {
		log.Println("QueryAllUser ", err)
		return err
	}
	for _, v := range result.Data.Records {
		avUsers[v.Uid] = v.Enable
	}
	return nil
}
func OnUserSub(lastBook *zhwatch.BookInfo, alreadySend, avUsers map[string]bool, recvUsrChan <-chan *wxpush.RecvUserData, client RedisConn) {
	for {
		select {
		case r := <-recvUsrChan:
			if _, ok := avUsers[r.GetUid()]; !ok {
				avUsers[r.GetUid()] = true
			}
			if _, ok := alreadySend[r.GetUid()]; !ok {
				SendBookInfo(lastBook, []string{r.GetUid()}, alreadySend, client)
			}
		}
	}
}
