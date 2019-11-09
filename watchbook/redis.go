package main

import (
	"encoding/json"
	"log"
	"mall/zhwatch"
	"time"

	"github.com/go-redis/redis"
)

type RedisConn struct {
	client *redis.Client
}

func NewRedisConn(ip, port string) RedisConn {
	var client = redis.NewClient(&redis.Options{
		Addr:     ip + ":" + port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return RedisConn{client}
}

//>func (c Pipeline) SScan(key string, cursor uint64, match string, count int64) *ScanCmd
func (c *RedisConn) GetAllAlready(key string, alreadySend map[string]bool) {
	if c.client == nil || alreadySend == nil {
		return
	}
	//Result() (keys []string, cursor uint64, err error)
	keys, cursor, err := c.client.SScan(key, 0, "", 100).Result()
	if err != nil {
		return
	}
	for _, k := range keys {
		alreadySend[k] = true
	}
	for cursor > 0 {
		keys, cursor, err = c.client.SScan(key, cursor, "", 100).Result()
		if err != nil {
			return
		}
		for _, k := range keys {
			alreadySend[k] = true
		}
	}
}

//func (c Client) SAdd(key string, members ...interface{}) *IntCmd
func (c *RedisConn) SetAlready(key, uid string) {
	if c.client == nil {
		return
	}
	i, err := c.client.SAdd(key, uid).Result()
	log.Println("SAdd", key, uid, i, err)
}

//func (c ClusterClient) Set(key string, value interface{}, expiration time.Duration) *StatusCmd
func (c *RedisConn) SetLastBook(key string, book *zhwatch.BookInfo) {
	if c.client == nil || book == nil {
		return
	}
	b, err := json.Marshal(book)
	if err != nil {
		log.Println("Marshal book err", err)
	}
	s, err := c.client.Set(key, string(b), 0*time.Second).Result()
	log.Println(s, err)
}
func (c *RedisConn) GetLastBook(key string, book *zhwatch.BookInfo) {
	if c.client == nil || book == nil {
		return
	}
	//func (c Client) Get(key string) *StringCmd
	b, err := c.client.Get(key).Result()
	if err != nil {
		return
	}
	json.Unmarshal([]byte(b), book)
	log.Printf("Get key %s %s\n", key, b)
}

//	func (c Client) Del(keys ...string) *IntCmd
func (c *RedisConn) Del(key string) {
	if c.client == nil {
		return
	}
	c.client.Del(key)
	log.Println("Del Key ", key)
}
