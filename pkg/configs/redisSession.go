package configs

import (
	"github.com/gin-contrib/sessions/redis"
)

func NewSessionStore() redis.Store {
	store, err := redis.NewStore(10, "tcp", "127.0.0.1:6379", "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", []byte("secret"))
	if err != nil {
		panic(err)
	}
	return store
}
