package configs

import (
	"fmt"

	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr   string
	DBName int
	Pass   string
	User   string
}

func NewTokenStore(conf *RedisConfig) *oredis.TokenStore {
	fmt.Println(conf.User)
	ts := oredis.NewRedisStore(&redis.Options{
		Addr:     conf.Addr,
		DB:       conf.DBName,
		Password: conf.Pass,
		Username: conf.User,
	})

	return ts
}
