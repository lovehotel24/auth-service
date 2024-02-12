package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/lovehotel24/auth-service/pkg/configs"
	"github.com/lovehotel24/auth-service/pkg/models"
)

const userKey = "userId"

func GetUsers(c *gin.Context) {
	var users []models.User
	configs.DB.Find(&users)
	c.JSON(http.StatusOK, &users)
}

func getUserById(userID interface{}) models.User {
	var user models.User
	configs.DB.Where("id = ?", userID).First(&user)
	return user
}

func getUserByPhone(phone string) models.User {
	var user models.User
	configs.DB.Where("phone = ?", phone).First(&user)
	return user
}

func CurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	user := getUserById(userID)
	c.JSON(http.StatusOK, &user)
}

func GetUser(c *gin.Context) {
	var user models.User
	configs.DB.Where("id = ?", c.Param("id")).First(&user)
	c.JSON(http.StatusOK, &user)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if hash, ok := generateHashPasswd(c, user.Password); ok {
		user.PasswordHash = hash
	} else {
		return
	}

	user.Password = ""
	configs.DB.Create(&user)
	c.JSON(http.StatusOK, &user)
}

func UpdateUser(c *gin.Context) {
	var updateUser models.User
	var currentUser models.User
	session := sessions.Default(c)
	userId := session.Get(userKey)
	fmt.Println("userID from Session of UpdateUser -> ", userId)
	user := getUserById(userId)
	if err := c.BindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	configs.DB.Model(models.User{}).Where("id = ?", user.Id).First(&currentUser)
	if currentUser.Id != user.Id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID for the update operation"})
		return
	}

	if hash, ok := generateHashPasswd(c, updateUser.Password); ok {
		user.PasswordHash = hash
	} else {
		return
	}

	if updateUser.Role != "" && updateUser.Role != user.Role && user.Role == "ADMIN" {
		user.Role = updateUser.Role
	}
	fmt.Println(updateUser.Phone)

	user.Phone = updateUser.Phone

	user.Name = updateUser.Name
	user.Password = ""

	if err := configs.DB.Save(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	configs.DB.Where("id = ?", c.Param("id")).Delete(&user)
	c.JSON(http.StatusOK, &user)
}

type forgetPass struct {
	Phone string `json:"phone"`
}

func ForgetPass(c *gin.Context) {
	var forget forgetPass
	if err := c.BindJSON(&forget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := getUserByPhone(forget.Phone)
	newReset := models.ResetPass{
		VerifyCode: configs.EncodeToString(6),
		UserId:     user.Id,
	}
	configs.DB.Create(&newReset)
	c.JSON(http.StatusOK, &newReset)
}

type resetPass struct {
	Code            string `json:"code"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func ResetPass(c *gin.Context) {
	var reset resetPass
	var forget models.ResetPass
	var user models.User
	if err := c.BindJSON(&reset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(reset.ConfirmPassword)
	if reset.Password != reset.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password and confirm password must be identical."})
		return
	}
	configs.DB.Model(&models.ResetPass{}).Where("verify_code = ?", reset.Code).First(&forget)
	configs.DB.Model(&models.User{}).Where("id = ?", forget.UserId).First(&user)

	if hash, ok := generateHashPasswd(c, reset.Password); ok {
		user.PasswordHash = hash
	} else {
		return
	}

	configs.DB.Save(&user)
	configs.DB.Model(&models.ResetPass{}).Delete(&forget)
	c.JSON(http.StatusOK, &user)
}

func generateHashPasswd(c *gin.Context, pass string) ([]byte, bool) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}
	return hash, true
}

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		user := getUserById(userId)

		if user.Role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
