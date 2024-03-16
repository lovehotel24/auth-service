package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func NewOauth2(db *gorm.DB, ts *oredis.TokenStore, clientStore oauth2.ClientStore) *server.Server {

	cfg := &manage.Config{
		AccessTokenExp:    time.Hour * 2,
		RefreshTokenExp:   time.Hour * 24 * 7,
		IsGenerateRefresh: true,
	}

	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetPasswordTokenCfg(cfg)

	manager.MapTokenStorage(ts)
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("jwt", []byte("secret"), jwt.SigningMethodHS512))
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

func (a API) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() == "/v1/user/login" || c.FullPath() == "/v1/user/register" || c.FullPath() == "/v1/user/forget_pass" || c.FullPath() == "/v1/user/reset_pass" {
			c.Next()
			return
		}

		token, err := a.SRV.ValidationBearerToken(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set(userKey, token.GetUserID())
		c.Next()
	}
}
