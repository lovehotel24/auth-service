package models

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id           uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4();"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone" gorm:"uniqueIndex"`
	Role         string    `json:"role"`
	Password     string    `json:"password"`
	PasswordHash []byte    `json:"-"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Id, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	return nil
}

type ResetPass struct {
	gorm.Model
	VerifyCode string
	UserId     uuid.UUID `gorm:"type:uuid;index"`
}

func NewAdmin() *User {
	admin := &User{
		Name:  "admin",
		Phone: "0634349640",
		Role:  "ADMIN",
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("hell123"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	admin.PasswordHash = hash

	return admin
}

func NewUser() *User {
	user := &User{
		Name:  "tester",
		Phone: "0634349641",
		Role:  "USER",
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("hell123"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	user.PasswordHash = hash

	return user
}
