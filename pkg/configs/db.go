package configs

import (
	"fmt"

	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lovehotel24/auth-service/pkg/models"
)

var DB *gorm.DB

type DBConfig struct {
	Host       string
	Port       string
	User       string
	Pass       string
	DBName     string
	SSLMode    string
	AdminPhone string
	AdminPass  string
	UserPhone  string
	UserPass   string
}

func Connect(conf *DBConfig, grpcClient userpb.UserServiceClient) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok", conf.Host, conf.User, conf.Pass, conf.DBName, conf.Port, conf.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	err = db.AutoMigrate(&models.DBUser{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.ResetPass{})
	if err != nil {
		panic(err)
	}
	admin := models.NewAdmin(conf.AdminPhone, conf.AdminPass, grpcClient)
	db.Create(&admin)
	tester := models.NewUser(conf.UserPhone, conf.UserPass, grpcClient)
	db.Create(&tester)
	DB = db
}
