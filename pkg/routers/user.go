package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"

	"github.com/lovehotel24/auth-service/pkg/controller"
)

func UserRouter(router *gin.Engine, srv *server.Server) {
	router.Use(ValidateToken(srv))
	v1UserRouter := router.Group("/v1")
	v1UserRouter.GET("/users", controller.GetUsers)
	v1UserRouter.GET("/user/:id", controller.GetUser)
	v1UserRouter.GET("/current_user", controller.CurrentUser)
	v1UserRouter.POST("/register", controller.CreateUser)
	v1UserRouter.DELETE("/user/:id", controller.DeleteUser)
	v1UserRouter.PUT("/user/:id", controller.UpdateUser)
	v1UserRouter.POST("/login", controller.Login())
}

func ValidateToken(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() == "/v1/login" || c.FullPath() == "/oauth/token" || c.FullPath() == "/v1/register" {
			c.Next()
			return
		}

		token, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("userID", token.GetUserID())
		c.Next()
	}
}
