package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id           uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4();"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
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
