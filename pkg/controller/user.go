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
	Password string `json:"password,omitempty"`
}

func NewUserId() uuid.UUID {
	userId, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("fail to create uuid")
	}
	return userId
}

func (a API) getDBUserById(userId interface{}) models.DBUser {
	var user models.DBUser
	a.DB.Where("id = ?", userId).First(&user)
	return user
}

func (a API) getDBUserByPhone(phone string) models.DBUser {
	var user models.DBUser
	a.DB.Where("phone = ?", phone).First(&user)
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

func (a API) CurrentUser(c *gin.Context) {
	userId, ok := c.Get(userKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to current user"})
		return
	}
	user, done := getBookUserById(userId.(string), a.Grpc)
	if !done {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func (a API) GetUser(c *gin.Context) {

	userId := c.Param("id")
	user, done := getBookUserById(userId, a.Grpc)
	if !done {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, &user)
}

func (a API) GetUsers(c *gin.Context) {
	limit := 10
	offset, _ := strconv.Atoi(c.Query("offset"))

	req := &userpb.GetAllUserRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	allUsers, err := a.Grpc.GetAllUsers(context.Background(), req)
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

func (a API) CreateUser(c *gin.Context) {
	var user User
	var dbUser models.DBUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := NewUserId()
	dbUser.Id = userId
	user.Id = userId.String()

	dbUser.Phone = user.Phone

	if hash, ok := generateHashPasswd(c, user.Password); ok {
		dbUser.PasswordHash = hash
	} else {
		return
	}
	user.Password = ""

	_, err := a.Grpc.CreateUser(context.Background(), &userpb.CreateUserRequest{User: &userpb.User{
		Id:    &userpb.UUID{Value: userId.String()},
		Name:  user.Name,
		Phone: user.Phone,
		Role:  user.Role,
	}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	a.DB.Create(&dbUser)
	c.JSON(http.StatusOK, &user)

}

func (a API) UpdateUser(c *gin.Context) {
	var (
		updateUser     User
		saveToDBUser   models.DBUser
		updateBookUser *userpb.User
	)

	userId, ok := c.Get(userKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	if err := c.BindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty param"})
		return
	}

	updateUser.Id = c.Param("id")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	saveToDBUser.Id = id
	updateBookUser = &userpb.User{Id: &userpb.UUID{Value: c.Param("id")}}

	currentUser, _ := getBookUserById(userId.(string), a.Grpc)
	bookUser, _ := getBookUserById(c.Param("id"), a.Grpc)
	dbUser := a.getDBUserById(c.Param("id"))

	switch currentUser.Role {
	case "ADMIN", "USER":
		if currentUser.Role == "USER" && c.Param("id") != userId.(string) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are not authorized to update other users' information"})
			return
		}

		if updateUser.Phone != "" {
			saveToDBUser.Phone = updateUser.Phone
			updateBookUser.Phone = updateUser.Phone
		} else {
			saveToDBUser.Phone = dbUser.Phone
			updateUser.Phone = dbUser.Phone
		}

		if updateUser.Password != "" {
			if hash, ok := generateHashPasswd(c, updateUser.Password); ok {
				saveToDBUser.PasswordHash = hash
			} else {
				return
			}
		} else {
			saveToDBUser.PasswordHash = dbUser.PasswordHash
		}

		if updateUser.Name != "" {
			updateBookUser.Name = updateUser.Name
		} else {
			updateUser.Name = bookUser.Name
		}

		if updateUser.Role != "" {
			updateBookUser.Role = updateUser.Role
		} else {
			updateUser.Role = bookUser.Role
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized user role for this operation"})
		return
	}

	if _, err = a.Grpc.UpdateUser(context.Background(), &userpb.UpdateUserRequest{User: updateBookUser}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform user update"})
		return
	}

	updateUser.Password = ""

	if err := a.DB.Save(saveToDBUser).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, &updateUser)
}

func (a API) DeleteUser(c *gin.Context) {
	var user models.DBUser
	if _, err := a.Grpc.DeleteUser(context.Background(), &userpb.DeleteUserRequest{Id: &userpb.UUID{Value: c.Param("id")}}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the user"})
		return
	}

	if err := a.DB.Where("id = ?", c.Param("id")).Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the user"})
		return
	}

	c.JSON(http.StatusOK, &user)

}

type forgetPass struct {
	Phone string `json:"phone"`
}

func (a API) ForgetPass(c *gin.Context) {
	var forget forgetPass
	if err := c.BindJSON(&forget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := a.getDBUserByPhone(forget.Phone)
	newReset := models.ResetPass{
		VerifyCode: configs.EncodeToString(6),
		UserId:     user.Id,
	}
	if err := a.DB.Create(&newReset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &newReset)
}

type resetPass struct {
	Code            string `json:"code"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (a API) ResetPass(c *gin.Context) {
	var reset resetPass
	var forget models.ResetPass
	var user models.DBUser
	if err := c.BindJSON(&reset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reset.Password != reset.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password and confirm password must be identical."})
		return
	}
	a.DB.Model(&models.ResetPass{}).Where("verify_code = ?", reset.Code).First(&forget)
	a.DB.Model(&models.DBUser{}).Where("id = ?", forget.UserId).First(&user)

	if hash, ok := generateHashPasswd(c, reset.Password); ok {
		user.PasswordHash = hash
	} else {
		return
	}

	a.DB.Save(&user)
	a.DB.Model(&models.ResetPass{}).Delete(&forget)
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

func (a API) OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get(userKey)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		user, done := getBookUserById(userId.(string), a.Grpc)
		if !done {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
			return
		}

		if user.Role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
