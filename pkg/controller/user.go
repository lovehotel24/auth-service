package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/lovehotel24/auth-service/pkg/configs"
	"github.com/lovehotel24/auth-service/pkg/models"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	configs.DB.Find(&users)
	c.JSON(http.StatusOK, &users)
}

func CurrentUser(c *gin.Context) {
	var user models.User
	userID, _ := c.Get("userID")
	configs.DB.Where("id = ?", userID).First(&user)
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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.PasswordHash = hash
	user.Password = ""
	configs.DB.Create(&user)
	c.JSON(200, &user)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	configs.DB.Where("id = ?", c.Param("id")).First(&user)
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	configs.DB.Save(&user)
	c.JSON(http.StatusOK, &user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	configs.DB.Where("id = ?", c.Param("id")).Delete(&user)
	c.JSON(http.StatusOK, &user)
}
