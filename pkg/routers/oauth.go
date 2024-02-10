package routers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	oauthErrs "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/server"
)

func OauthRouter(router *gin.Engine, srv *server.Server) {
	router.POST("/oauth/token", func(c *gin.Context) {
		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			return
		}
	})

	router.GET("/oauth/validate_token", func(c *gin.Context) {
		if token, err := srv.ValidationBearerToken(c.Request); err != nil {
			res := gin.H{"err_desc": err.Error()}
			switch {
			case errors.Is(err, oauthErrs.ErrInvalidAccessToken):
				res["err_no"] = "-1001"
			case errors.Is(err, oauthErrs.ErrExpiredAccessToken):
				res["err_no"] = "-1002"
			case errors.Is(err, oauthErrs.ErrExpiredRefreshToken):
				res["err_no"] = "-1003"
			default:
				res["err_no"] = "-1000"
			}
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
				"user_id":    token.GetUserID(),
			})
		}
	})
}
