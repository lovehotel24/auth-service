package pkg

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/lovehotel24/auth-service/pkg/configs"
	"github.com/lovehotel24/auth-service/pkg/controller"
	"github.com/lovehotel24/auth-service/pkg/routers"
)

func Run() {
	router := gin.New()
	configs.Connect()
	tokenStore := configs.NewTokenStore()
	sessionStore := configs.NewSessionStore()
	oauthSvr := controller.NewOauth2(configs.DB, tokenStore)
	//oauthSvr.SetPasswordAuthorizationHandler(controller.PasswordAuthorizationHandler(configs.DB))
	router.Use(gin.Logger(), sessions.Sessions("session", sessionStore))
	routers.UserRouter(router, oauthSvr)
	routers.OauthRouter(router, oauthSvr)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
