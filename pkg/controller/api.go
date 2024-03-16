package controller

import (
	"github.com/go-oauth2/oauth2/v4/server"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type API struct {
	DB          *gorm.DB
	Log         *logrus.Logger
	Grpc        userpb.UserServiceClient
	OauthConfig oauth2.Config
	TS          *oredis.TokenStore
	SRV         *server.Server
}

func NewApp(db *gorm.DB, log *logrus.Logger, grpc userpb.UserServiceClient, oConfig oauth2.Config, ts *oredis.TokenStore, srv *server.Server) *API {
	return &API{
		DB:          db,
		Log:         log,
		Grpc:        grpc,
		OauthConfig: oConfig,
		TS:          ts,
		SRV:         srv,
	}
}
