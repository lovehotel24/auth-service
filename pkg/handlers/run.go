package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/lovehotel24/auth-service/pkg/auth/oauth2"
	"github.com/lovehotel24/auth-service/pkg/foundation/validate"
	"github.com/lovehotel24/auth-service/pkg/model/user"
	"github.com/lovehotel24/auth-service/pkg/sys/database"
)

var build = "develop"

func Run(log *zap.SugaredLogger) error {
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	// =================================================================================================================
	// Database Support

	// Create connectivity to the database.
	log.Infow("startup", "status", "initializing database support", "host", "cfg.DB.Host")

	db, err := database.Open(database.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         "localhost",
		Name:         "users",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		DisableTLS:   true,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", "localhost")
		db.Close()
	}()

	// =================================================================================================================
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	// =================================================================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	userStore := user.NewStore(log, db)
	oas := oauth2.NewOauthServer()
	oas.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler(userStore))

	// Constructs the mux for the API calls.
	apiMux := APIMux(APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		DB:       db,
		OAS:      oas,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      apiMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =================================================================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func PasswordAuthorizationHandler(store user.Store) func(context.Context, string, string, string) (string, error) {
	return func(ctx context.Context, clientID, phone, password string) (string, error) {
		now := time.Now()

		if clientID == "222222" {
			userID, err := store.Authenticate(ctx, now, phone, password)
			if err != nil {
				switch validate.Cause(err) {
				case database.ErrNotFound:
					return "", validate.NewRequestError(err, http.StatusNotFound)
				case database.ErrAuthenticationFailure:
					return "", validate.NewRequestError(err, http.StatusUnauthorized)
				default:
					return "", fmt.Errorf("authenticating: %w", err)
				}
			}
			return userID, nil
		}
		return "", errors.ErrUnauthorizedClient
	}
}
