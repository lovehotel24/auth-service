package configs

import (
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
)

func NewTokenStore() *oredis.TokenStore {
	ts := oredis.NewRedisStore(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       15,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})

	return ts
}
