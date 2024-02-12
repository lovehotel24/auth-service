package routers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"

	"github.com/lovehotel24/auth-service/pkg/controller"
)

const userKey = "userId"

func UserRouter(router *gin.Engine, srv *server.Server) {
	router.Use(ValidateToken(srv))
	v1UserRouter := router.Group("/v1")
	v1UserRouter.GET("/users", controller.OnlyAdmin(), controller.GetUsers)
	v1UserRouter.GET("/user/:id", controller.GetUser)
	v1UserRouter.GET("/current_user", controller.CurrentUser)
	v1UserRouter.POST("/register", controller.CreateUser)
	v1UserRouter.DELETE("/user/:id", controller.OnlyAdmin(), controller.DeleteUser)
	v1UserRouter.PUT("/user/:id", controller.UpdateUser)
	v1UserRouter.POST("/forget_pass", controller.ForgetPass)
	v1UserRouter.POST("/reset_pass", controller.ResetPass)
	v1UserRouter.POST("/login", controller.Login())
}

func ValidateToken(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() == "/v1/login" || c.FullPath() == "/oauth/token" || c.FullPath() == "/v1/register" || c.FullPath() == "/v1/forget_pass" || c.FullPath() == "/v1/reset_pass" {
			c.Next()
			return
		}

		token, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		session := sessions.Default(c)
		session.Set(userKey, token.GetUserID()) // In real world usage you'd set this to the users ID
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		c.Set("userID", token.GetUserID())
		c.Next()
	}
}
