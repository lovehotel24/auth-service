package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"

	"github.com/lovehotel24/auth-service/pkg/auth/oauth2"
	"github.com/lovehotel24/auth-service/pkg/foundation/validate"
	"github.com/lovehotel24/auth-service/pkg/model/user"
	"github.com/lovehotel24/auth-service/pkg/sys/database"
)

type Handlers struct {
	Oauth2 *oauth2.OAuthServer
}

func (h Handlers) Token(w http.ResponseWriter, r *http.Request) {

	err := h.Oauth2.Srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {

	token, err := h.Oauth2.Srv.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

//func main() {
//	oauth := oauth2.NewOauthServer()
//
//	db, err := database.Open(database.Config{
//		User:         "postgres",
//		Password:     "postgres",
//		Host:         "localhost",
//		Name:         "users",
//		MaxIdleConns: 0,
//		MaxOpenConns: 0,
//		DisableTLS:   true,
//	})
//	if err != nil {
//		_ = fmt.Errorf("connecting to db: %w", err)
//	}
//	log, _ := logger.New("AUTH")
//	userStore := user.NewStore(log, db)
//	oauth.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler(userStore))
//	ouh := &Handlers{Oauth2: oauth}
//
//	http.HandleFunc("/v1/oauth/token", ouh.Token)
//	http.HandleFunc("/v1/oauth/test", ouh.Test)
//	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8081), nil))
//}

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
