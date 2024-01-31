package handlers

import (
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lovehotel24/auth-service/pkg/auth/oauth2"
	"github.com/lovehotel24/auth-service/pkg/foundation/web"
	"github.com/lovehotel24/auth-service/pkg/handlers/v1/oauth"
	"github.com/lovehotel24/auth-service/pkg/handlers/v1/testgrp"
	"github.com/lovehotel24/auth-service/pkg/handlers/v1/usergrp"
	"github.com/lovehotel24/auth-service/pkg/model/user"
	"github.com/lovehotel24/auth-service/pkg/sys/middleware"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	DB       *sqlx.DB
	OAS      *oauth2.OAuthServer
}

func APIMux(cfg APIMuxConfig) *web.App {

	// Construct the web.app which holds all routes.
	app := web.NewApp(
		cfg.Shutdown,
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
	)

	// Load the routes for the different versions of the API.
	v1(app, cfg)

	return app
}

// v1 binds all the version 1 routes.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, version, "/test", tgh.Test)

	ugh := usergrp.Handlers{
		User: user.NewStore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodGet, version, "/users/:page/:rows", ugh.Query)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID)
	app.Handle(http.MethodPost, version, "/users", ugh.Create)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete)
	app.Handle(http.MethodPost, version, "/users/login", ugh.Login)

	ogh := oauth.Handlers{
		Oauth2: cfg.OAS,
	}

	app.Handle(http.MethodPost, version, "/oauth/authorize", ogh.Authorize)
	app.Handle(http.MethodGet, version, "/oauth/token", ogh.Token)
	app.Handle(http.MethodGet, version, "/oauth/test", ogh.Test)

	app.Handle(http.MethodGet, version, "/", ogh.Index)
	app.Handle(http.MethodGet, version, "/oauth/oauth2", ogh.OAuth2)
	app.Handle(http.MethodGet, version, "/oauth/refresh", ogh.Refresh)
	app.Handle(http.MethodGet, version, "/oauth/try", ogh.Try)
	app.Handle(http.MethodGet, version, "/oauth/pwd", ogh.PWD)
	app.Handle(http.MethodGet, version, "/oauth/client", ogh.Client)
}
