package configs

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lovehotel24/auth-service/pkg/models"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.ResetPass{})
	if err != nil {
		panic(err)
	}
	admin := models.NewAdmin()
	db.Create(&admin)
	tester := models.NewUser()
	db.Create(&tester)
	DB = db
}
