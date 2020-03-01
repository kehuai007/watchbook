package server

import (
	"bytes"
	"fmt"
	"log"
	"misu/watchbook/lib"

	"github.com/kehuai007/wxpush"
)

type (
	Server struct {
		// appName-Request
		Request    []request
		ListenPort int
	}
	request struct {
		s           lib.WatchBookServer
		appToken    string
		useList     map[string]bool
		minInterval int
	}
)

func (r *request) QueryUser() {
	u, err := wxpush.QueryAllUser(r.appToken)
	if err == nil && u != nil {
		r.useList = make(map[string]bool)
		var buffer bytes.Buffer
		for _, usr := range u.Data.Records {
			r.useList[usr.Uid] = true
			buffer.WriteString(usr.Uid)
			buffer.WriteString(",")
		}
		log.Printf("%s aapToken %s useList(%d) %s\n", r.s.GetBookRequest().Name, r.appToken, len(r.useList), buffer.String())
	}
}
func NewServer(port int) *Server {
	return &Server{ListenPort: port}
}
func (s *Server) PushRequest(req lib.WatchBookServer, appToken string, minute int) {
	s.Request = append(s.Request, request{
		s:           req,
		appToken:    appToken,
		minInterval: minute,
	})
}
func (s *Server) SendMsg(book *lib.Book, r request) {
	if book != nil && len(r.useList) <= 0 {
		return
	}
	msg := wxpush.NewMessage(r.appToken)
	msg.SetContent(fmt.Sprintf("《%s》更新啦！最新 %s。", book.Name, book.Title))
	msg.SetUrl(book.Url)
	for u, _ := range r.useList {
		msg.AddUId(u)
	}
	wxpush.SendMessage(*msg)
}
func (s *Server) Run() {
	for _, r := range s.Request {
		r.QueryUser()
		go func(ch chan *lib.Book) {
			for {
				select {
				case book := <-ch:
					s.SendMsg(book, r)
				}
			}
		}(r.s.Run(r.minInterval))
	}
	usrCh := make(chan *wxpush.RecvUserData)
	go func() {
		for {
			select {
			case u := <-usrCh:
				for _, r := range s.Request {
					if _, ok := r.useList[u.GetUid()]; !ok {
						r.QueryUser()
					}
				}
			}
		}
	}()
	log.Printf("start listen %d", s.ListenPort)
	err := wxpush.ServerAndListen(s.ListenPort, usrCh)
	if err != nil {
		panic(err)
	}
}
