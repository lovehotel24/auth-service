package handlers

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/lovehotel24/auth-service/pkg/foundation/web"
	"github.com/lovehotel24/auth-service/pkg/handlers/v1/testgrp"
	"github.com/lovehotel24/auth-service/pkg/sys/middleware"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
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
}
