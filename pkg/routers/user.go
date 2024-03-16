package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/lovehotel24/auth-service/pkg/controller"
)

const (
	userKey = "userId"
)

func UserRouter(router *gin.Engine, api *controller.API) {
	v1Route := router.Group("/v1")
	userRouter := v1Route.Group("/user")
	userRouter.Use(api.ValidateToken())
	userRouter.GET("/", api.GetUsers)
	userRouter.GET("/:id", api.GetUser)
	userRouter.GET("/current_user", api.CurrentUser)
	userRouter.POST("/register", api.CreateUser)
	userRouter.DELETE("/:id", api.OnlyAdmin(), api.DeleteUser)
	userRouter.PUT("/:id", api.UpdateUser)
	userRouter.POST("/forget_pass", api.ForgetPass)
	userRouter.POST("/reset_pass", api.ResetPass)
	userRouter.GET("/logout", api.Logout())
	userRouter.POST("/login", api.Login())
}
