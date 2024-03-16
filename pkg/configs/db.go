package configs

import (
	"fmt"

	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/lovehotel24/auth-service/pkg/models"
)

type DBConfig struct {
	host     string
	port     string
	user     string
	pass     string
	name     string
	sslMode  string
	timeZone string
}

func (c DBConfig) WithHost(host string) DBConfig {
	c.host = host
	return c
}

func (c DBConfig) WithPort(port string) DBConfig {
	c.port = port
	return c
}

func (c DBConfig) WithUser(user string) DBConfig {
	c.user = user
	return c
}

func (c DBConfig) WithPass(pass string) DBConfig {
	c.pass = pass
	return c
}

func (c DBConfig) WithName(name string) DBConfig {
	c.name = name
	return c
}

func (c DBConfig) WithSecure(ssl bool) DBConfig {
	if ssl {
		c.sslMode = "enable"
	}
	return c
}

func (c DBConfig) WithTZ(tz string) DBConfig {
	c.timeZone = tz
	return c
}

func NewDB(conf DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", conf.host, conf.user, conf.pass, conf.name, conf.port, conf.sslMode, conf.timeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	return db, nil
}

func NewDBConfig() DBConfig {
	return DBConfig{
		host:     "localhost",
		port:     "5432",
		user:     "postgres",
		pass:     "postgres",
		name:     "postgres",
		sslMode:  "disable",
		timeZone: "Asia/Bangkok",
	}
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.DBUser{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.ResetPass{})
	if err != nil {
		return err
	}
	return err
}

type DefaultUser struct {
	admPh    string
	admPass  string
	userPh   string
	UserPass string
}

func NewDefaultUser() DefaultUser {
	return DefaultUser{
		admPh:    "0612345678",
		admPass:  "topSecret",
		userPh:   "0601234567",
		UserPass: "lowSecret",
	}
}

func (d DefaultUser) WithDefaultAdminPhone(phone string) DefaultUser {
	d.admPh = phone
	return d
}

func (d DefaultUser) WithDefaultAdminPass(pass string) DefaultUser {
	d.admPh = pass
	return d
}

func (d DefaultUser) WithDefaultUserPhone(phone string) DefaultUser {
	d.userPh = phone
	return d
}

func (d DefaultUser) WithDefaultUserPass(pass string) DefaultUser {
	d.UserPass = pass
	return d
}

func Seed(db *gorm.DB, defaultUser DefaultUser, grpcClient userpb.UserServiceClient) error {
	admin := models.NewAdmin(defaultUser.admPh, defaultUser.admPass, grpcClient)
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	tester := models.NewUser(defaultUser.userPh, defaultUser.UserPass, grpcClient)
	if err := db.Create(&tester).Error; err != nil {
		return err
	}

	return nil
}
