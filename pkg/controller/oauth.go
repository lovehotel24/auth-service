package controller

import (
	"log"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func NewOauth2(db *gorm.DB) *server.Server {

	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:8080",
	})

	cfg := &manage.Config{
		AccessTokenExp:    time.Hour * 2,
		RefreshTokenExp:   time.Hour * 24 * 7,
		IsGenerateRefresh: true,
	}

	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetPasswordTokenCfg(cfg)

	// token store
	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       15,
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	}))
	//manager.MapAccessGenerate(generates.NewJWTAccessGenerate("jwt", []byte("pibigstar"), jwt.SigningMethodHS512))
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	srv.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler(db))

	return srv
}
