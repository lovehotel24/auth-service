package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/lovehotel24/auth-service/pkg/models"
)

type userLogin struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func PasswordAuthorizationHandler(db *gorm.DB) func(context.Context, string, string, string) (string, error) {
	return func(ctx context.Context, clientID, phone, password string) (string, error) {
		var user models.DBUser
		if clientID != viper.GetString("client-id") {
			return "", errors.ErrUnauthorizedClient
		}
		err := db.Model(&user).Where("phone = ?", phone).First(&user).Error
		if err != nil {
			return "", errors.ErrUnauthorizedClient
		}
		err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
		if err != nil {
			return "", errors.ErrUnauthorizedClient
		}
		return user.Id.String(), nil
	}
}

func Login(config oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user userLogin
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := config.PasswordCredentialsToken(context.Background(), user.Phone, user.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, token)
	}
}

func Logout(ts *oredis.TokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := BearerAuth(c.Request)
		if !ok {
			c.Redirect(http.StatusPermanentRedirect, "/")
			return
		}
		err := ts.RemoveByAccess(context.Background(), token)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "DBUser Sign out successfully",
		})
	}
}

func BearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access_token")
	}

	return token, token != ""
}
