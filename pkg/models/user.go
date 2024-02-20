package models

import (
	"fmt"

	"github.com/google/uuid"
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

func NewAdmin(phone, pass string) *DBUser {
	admin := &DBUser{
		Phone: phone,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	admin.PasswordHash = hash

	return admin
}

func NewUser(phone, pass string) *DBUser {
	user := &DBUser{
		Phone: phone,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	user.PasswordHash = hash

	return user
}
