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
	v1Route := router.Group("/v1")
	userRouter := v1Route.Group("/user")
	userRouter.Use(ValidateToken(srv))
	userRouter.GET("/", controller.GetUsers(grpcClient))
	userRouter.GET("/:id", controller.GetUser(grpcClient))
	userRouter.GET("/current_user", controller.CurrentUser(grpcClient))
	userRouter.POST("/register", controller.CreateUser(grpcClient))
	userRouter.DELETE("/:id", controller.OnlyAdmin(grpcClient), controller.DeleteUser(grpcClient))
	userRouter.PUT("/:id", controller.UpdateUser(grpcClient))
	userRouter.POST("/forget_pass", controller.ForgetPass)
	userRouter.POST("/reset_pass", controller.ResetPass)
	userRouter.GET("/logout", controller.Logout(ts))
	userRouter.POST("/login", controller.Login(client))
}

func ValidateToken(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() == "/v1/user/login" || c.FullPath() == "/v1/user/register" || c.FullPath() == "/v1/user/forget_pass" || c.FullPath() == "/v1/user/reset_pass" {
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
