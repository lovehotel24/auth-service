package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"golang.org/x/oauth2"

	"github.com/lovehotel24/auth-service/pkg/controller"
)

const (
	userKey = "userId"
)

func UserRouter(router *gin.Engine, srv *server.Server, ts *oredis.TokenStore, client oauth2.Config, grpcClient userpb.UserServiceClient) {
	v1UserRouter := router.Group("/v1")
	v1UserRouter.Use(ValidateToken(srv))
	v1UserRouter.GET("/users", controller.GetUsers(grpcClient))
	v1UserRouter.GET("/user/:id", controller.GetUser(grpcClient))
	v1UserRouter.GET("/current_user", controller.CurrentUser(grpcClient))
	v1UserRouter.POST("/register", controller.CreateUser(grpcClient))
	v1UserRouter.DELETE("/user/:id", controller.OnlyAdmin(grpcClient), controller.DeleteUser(grpcClient))
	v1UserRouter.PUT("/user/:id", controller.UpdateUser(grpcClient))
	v1UserRouter.POST("/forget_pass", controller.ForgetPass)
	v1UserRouter.POST("/reset_pass", controller.ResetPass)
	v1UserRouter.GET("/hello", controller.Hello)
	v1UserRouter.GET("/logout", controller.Logout(ts))
	v1UserRouter.POST("/login", controller.Login(client))
}

func ValidateToken(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() == "/v1/login" || c.FullPath() == "/v1/register" || c.FullPath() == "/v1/forget_pass" || c.FullPath() == "/v1/reset_pass" {
			c.Next()
			return
		}

		token, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set(userKey, token.GetUserID())
		c.Next()
	}
}
