package models

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DBUser struct {
	gorm.Model
	Id           uuid.UUID `gorm:"primary_key;type:uuid;"`
	Phone        string    `gorm:"<-:create;uniqueIndex"`
	PasswordHash []byte
}

//
//func (user *DBUser) BeforeCreate(tx *gorm.DB) (err error) {
//	user.Id, err = uuid.NewUUID()
//	if err != nil {
//		return err
//	}
//	return nil
//}

type ResetPass struct {
	gorm.Model
	VerifyCode string
	UserId     uuid.UUID `gorm:"type:uuid;index"`
}

func NewAdmin(phone, pass string, grpcClient userpb.UserServiceClient) *DBUser {
	userId := NewUserId()
	_, err := grpcClient.CreateUser(context.Background(), &userpb.CreateUserRequest{User: &userpb.User{
		Id:    &userpb.UUID{Value: userId.String()},
		Name:  "admin",
		Phone: phone,
		Role:  "ADMIN",
	}})
	if err != nil {
		fmt.Printf("Fail to create admin user: %v", err.Error())
	}
	admin := &DBUser{
		Id:    userId,
		Phone: phone,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	admin.PasswordHash = hash

	return admin
}

func NewUser(phone, pass string, grpcClient userpb.UserServiceClient) *DBUser {
	userId := NewUserId()
	_, err := grpcClient.CreateUser(context.Background(), &userpb.CreateUserRequest{User: &userpb.User{
		Id:    &userpb.UUID{Value: userId.String()},
		Name:  "test",
		Phone: phone,
		Role:  "USER",
	}})
	if err != nil {
		fmt.Printf("Fail to create test user: %v", err.Error())
	}
	user := &DBUser{
		Id:    userId,
		Phone: phone,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	user.PasswordHash = hash

	return user
}

func NewUserId() uuid.UUID {
	userId, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("fail to create uuid")
	}
	return userId
}
