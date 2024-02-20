package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"golang.org/x/crypto/bcrypt"

	"github.com/lovehotel24/auth-service/pkg/configs"
	"github.com/lovehotel24/auth-service/pkg/models"
)

const (
	userKey = "userId"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

func NewUserId() uuid.UUID {
	userId, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("fail to create uuid")
	}
	return userId
}

func getDBUserById(userID interface{}) models.DBUser {
	var user models.DBUser
	configs.DB.Where("id = ?", userID).First(&user)
	return user
}

func getDBUserByPhone(phone string) models.DBUser {
	var user models.DBUser
	configs.DB.Where("phone = ?", phone).First(&user)
	return user
}

func getBookUserById(userId string, grpcClient userpb.UserServiceClient) (User, bool) {
	getUser, err := grpcClient.GetUser(context.Background(), &userpb.GetUserRequest{Id: &userpb.UUID{Value: userId}})
	if err != nil {
		return User{}, false
	}
	gUser := getUser.GetUser()
	user := User{
		Id:    gUser.GetId().GetValue(),
		Name:  gUser.GetName(),
		Phone: gUser.GetPhone(),
		Role:  gUser.GetRole(),
	}
	return user, true
}

func CurrentUser(grpcClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(userKey)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to current user"})
			return
		}
		user, done := getBookUserById(userID.(string), grpcClient)
		if !done {
			return
		}
		c.JSON(http.StatusOK, &user)
	}
}

func GetUser(grpcClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")
		user, done := getBookUserById(userId, grpcClient)
		if !done {
			return
		}
		configs.DB.Where("id = ?").First(&user)
		c.JSON(http.StatusOK, &user)
	}
}

func GetUsers(grpcClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		//var users []User
		//blank := &empty.Empty{}
		//allUsers, err := grpcClient.GetAllUsers(context.Background(), blank)
		//if err != nil {
		//	return
		//}
		//for _, v := range allUsers.GetUsers() {
		//	user := User{
		//		Id:    v.GetId().GetValue(),
		//		Name:  v.GetName(),
		//		Phone: v.GetPhone(),
		//		Role:  v.GetRole(),
		//	}
		//	users = append(users, user)
		//}
		//c.JSON(http.StatusOK, &users)
		limit := 10
		offset, _ := strconv.Atoi(c.Query("offset"))

		req := &userpb.GetAllUserRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		allUsers, err := grpcClient.GetAllUsers(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		var users []User
		for _, v := range allUsers.GetUsers() {
			user := User{
				Id:    v.GetId().GetValue(),
				Name:  v.GetName(),
				Phone: v.GetPhone(),
				Role:  v.GetRole(),
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, &users)
	}
}

func CreateUser(grpcClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		var dbUser models.DBUser
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId := NewUserId()
		dbUser.Id = userId

		dbUser.Phone = user.Phone

		if hash, ok := generateHashPasswd(c, user.Password); ok {
			dbUser.PasswordHash = hash
		} else {
			return
		}
		user.Password = ""

		configs.DB.Create(&dbUser)
		_, err := grpcClient.CreateUser(context.Background(), &userpb.CreateUserRequest{User: &userpb.User{
			Id:    &userpb.UUID{Value: userId.String()},
			Name:  user.Name,
			Phone: user.Phone,
			Role:  user.Role,
		}})
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, &user)
	}
}

func UpdateUser(grpcClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateUser User
		var dbUpdateUser models.DBUser
		userId, ok := c.Get(userKey)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		user := getDBUserById(userId)
		if err := c.BindJSON(&updateUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//configs.DB.Model(models.DBUser{}).Where("id = ?", user.Id).First(&currentUser)
		if updateUser.Id != userId.(string) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID for the update operation"})
			return
		}
		updateBookUser := &userpb.User{
			Id: &userpb.UUID{Value: userId.(string)},
		}

		if updateUser.Phone != "" {
			dbUpdateUser.Phone = updateUser.Phone
			updateBookUser.Phone = updateUser.Phone
		} else {
			dbUpdateUser.Phone = user.Phone
		}

		if updateUser.Password != "" {
			if hash, ok := generateHashPasswd(c, updateUser.Password); ok {
				dbUpdateUser.PasswordHash = hash
			} else {
				return
			}
		} else {
			dbUpdateUser.PasswordHash = user.PasswordHash
		}

		if updateUser.Name != "" {
			updateBookUser.Name = updateUser.Name
		}

		if updateUser.Role != "" {
			updateBookUser.Role = updateUser.Role
		}

		_, err := grpcClient.UpdateUser(context.Background(), &userpb.UpdateUserRequest{User: updateBookUser})
		if err != nil {
			return
		}

		updateUser.Password = ""

		if err := configs.DB.Save(dbUpdateUser).Error; err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		c.JSON(http.StatusOK, &user)
	}
}

func DeleteUser(c *gin.Context) {
	var user models.DBUser
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
	user := getDBUserByPhone(forget.Phone)
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
	var user models.DBUser
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
	configs.DB.Model(&models.DBUser{}).Where("id = ?", forget.UserId).First(&user)

	if hash, ok := generateHashPasswd(c, reset.Password); ok {
		user.PasswordHash = hash
	} else {
		return
	}

	configs.DB.Save(&user)
	configs.DB.Model(&models.ResetPass{}).Delete(&forget)
	c.JSON(http.StatusOK, &user)
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
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
		//userId, ok := c.Get(userKey)
		//if !ok {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		//	c.Abort()
		//	return
		//}
		//user := getDBUserById(userId)
		//
		//if user.Role != "ADMIN" {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		//	c.Abort()
		//	return
		//}

		c.Next()
	}
}
