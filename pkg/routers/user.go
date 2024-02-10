package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/lovehotel24/auth-service/pkg/controller"
)

func UserRouter(router *gin.Engine) {
	//router.Use(minddleware.Auth)
	router.GET("/", controller.GetUsers)
	router.POST("/", controller.CreateUser)
	router.DELETE("/:id", controller.DeleteUser)
	router.PUT("/:id", controller.UpdateUser)
	router.POST("/login", controller.Login())
}
