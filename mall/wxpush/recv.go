package wxpush

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RecvUserData struct {
	Action string `json:"action"`
	Data   struct {
		AppKey  string `json:"app_key"`
		AppName string `json:"app_name"`
		Source  string `json:"source"`
		Time    int64  `json:"time"`
		Uid     string `json:"uid"`
		Extra   string `json:"extra"`
	} `json:"data"`
}

func (r *RecvUserData) GetUid() string {
	return r.Data.Uid
}

func ServerAndListen(port int, c chan<- *RecvUserData) error {

	recvHandle := func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			w.Write([]byte("Server is Running"))
			w.WriteHeader(http.StatusOK)
		case http.MethodPost:
			s, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Println("read body err ", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r := RecvUserData{}
			err = json.Unmarshal(s, &r)
			if err != nil {
				log.Println("Unmarshal err ", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			go func() {
				c <- &r
			}()
			w.WriteHeader(http.StatusOK)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	http.HandleFunc("/", recvHandle)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
