package configs

import (
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr   string
	DBName int
	Pass   string
}

func NewTokenStore(conf *RedisConfig) *oredis.TokenStore {
	ts := oredis.NewRedisStore(&redis.Options{
		Addr:     conf.Addr,
		DB:       conf.DBName,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})

	return ts
}
