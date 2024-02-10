package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/lovehotel24/auth-service/pkg/models"
)

const (
	authServerURL = "http://localhost:8080"
)

var (
	config = oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalToken *oauth2.Token
)

type userLogin struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func PasswordAuthorizationHandler(db *gorm.DB) func(context.Context, string, string, string) (string, error) {
	return func(ctx context.Context, clientID, phone, password string) (string, error) {
		var user models.User
		if clientID != "222222" {
			return "", errors.ErrUnauthorizedClient
		}
		err := db.Model(&user).Where("phone = ?", phone).First(&user).Error
		if err != nil {
			return "", errors.ErrUnauthorizedClient
		}
		fmt.Println(user.Password)
		err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
		if err != nil {
			return "", errors.ErrUnauthorizedClient
		}
		return user.Id.String(), nil
	}
}

func Login() gin.HandlerFunc {
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
		globalToken = token
		c.JSON(200, token)
	}
}
