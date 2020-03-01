package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"misu/watchbook/lib"
	"misu/watchbook/server"
	"os"
)

// 配置
type Config struct {
	Type string `json:"type"`
	BookUrl string `json:"book_url"`
	Name string `json:"name"`
	BookID int `json:"book_id"`
	PushToken string `json:"push_token"`
	LoopPeriod int `json:"loop_period"`
}

func LoadConfigData()  (map[string][]Config, error){
	filePtr, err := os.Open("watchbook.json")
	if err != nil {
		fmt.Println("Open file failed [Err:%s]", err.Error())
		return nil,err
	}
	defer filePtr.Close()
	b ,err :=ioutil.ReadAll(filePtr)
	if err != nil {
		fmt.Println("ReadAll file failed [Err:%s]", err.Error())
		return nil,err
	}
	var conf map[string][]Config
    err = json.Unmarshal(b,&conf)
	if err != nil {
		fmt.Println("Unmarshal file failed [Err:%s]", err.Error())
		return nil,err
	}
	return conf,nil
}

func main() {

	conf,err:=LoadConfigData()
	if err != nil{
		return
	}
	s := server.NewServer(9010)
	for name,c := range conf{
		if name =="qidian" {
			for _,d:=range c{
				r := lib.NewQiDian(lib.BookRequest{
					BookUrl: d.BookUrl,
					Name:    d.Name,
					BookId:  d.BookID,
				})
				s.PushRequest(r, d.PushToken,d.LoopPeriod)
			}
		}else if name == "zongheng"{
			for _,d:=range c{
				r := lib.NewZoneHen(lib.BookRequest{
					BookUrl: d.BookUrl,
					Name:    d.Name,
					BookId:  d.BookID,
				})
				s.PushRequest(r, d.PushToken,d.LoopPeriod)
			}
		}
	}
	s.Run()

}
