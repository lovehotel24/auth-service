package pkg

import (
	"github.com/gin-gonic/gin"

	"github.com/lovehotel24/auth-service/pkg/configs"
	"github.com/lovehotel24/auth-service/pkg/controller"
	"github.com/lovehotel24/auth-service/pkg/routers"
)

func Run() {
	router := gin.New()
	configs.Connect()
	oauthSvr := controller.NewOauth2(configs.DB)
	//oauthSvr.SetPasswordAuthorizationHandler(controller.PasswordAuthorizationHandler(configs.DB))
	router.Use(gin.Logger())
	routers.UserRouter(router, oauthSvr)
	routers.OauthRouter(router, oauthSvr)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
