package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/lovehotel24/auth-service/pkg/foundation/validate"
	"github.com/lovehotel24/auth-service/pkg/foundation/web"
)

const (
	authServerURL = "http://localhost:8080"
)

var (
	config = oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/v1/oauth/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/v1/oauth/authorize",
			TokenURL: authServerURL + "/v1/oauth/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

func (h Handlers) Index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("I'm Index ", r.RequestURI)
	u := config.AuthCodeURL("xyz",
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	fmt.Println(u)

	return web.Redirect(w, r, u, http.StatusFound)
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

func (h Handlers) OAuth2(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Println("I'm OAuth2 ", r.RequestURI)
	r.ParseForm()
	state := r.Form.Get("state")
	fmt.Println("I'm state ", state)
	if state != "xyz" {
		return validate.NewRequestError(fmt.Errorf("state invalid"), http.StatusBadRequest)
	}
	code := r.Form.Get("code")
	fmt.Println("I'm code ", code)
	if code == "" {
		return validate.NewRequestError(fmt.Errorf("code not found"), http.StatusBadRequest)
	}
	token, err := config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
	fmt.Println("I'm token ", token)
	if err != nil {
		return validate.NewRequestError(err, http.StatusInternalServerError)
	}
	globalToken = token

	return web.Respond(ctx, w, token, http.StatusOK)
}

func (h Handlers) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("I'm Refresh ", r.RequestURI)
	if globalToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	globalToken.Expiry = time.Now()
	token, err := config.TokenSource(context.Background(), globalToken).Token()
	if err != nil {
		return validate.NewRequestError(err, http.StatusInternalServerError)
	}

	globalToken = token

	return web.Respond(ctx, w, token, http.StatusOK)
}

func (h Handlers) Try(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("I'm Try ", r.RequestURI)
	if globalToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	resp, err := http.Get(fmt.Sprintf("%s/v/oauth/test?access_token=%s", authServerURL, globalToken.AccessToken))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return validate.NewRequestError(err, http.StatusBadRequest)
	}
	defer resp.Body.Close()

	//io.Copy(w, resp.Body)

	return web.Respond(ctx, w, resp.Body, http.StatusOK)
}

func (h Handlers) PWD(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Println("I'm PWD ", r.RequestURI)

	token, err := config.PasswordCredentialsToken(context.Background(), "0634349640", "aeiou1234")
	fmt.Println("I'm token ", token)
	if err != nil {
		return validate.NewRequestError(err, http.StatusInternalServerError)
	}

	globalToken = token

	return web.Respond(ctx, w, token, http.StatusOK)
}

func (h Handlers) Client(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("I'm Client ", r.RequestURI)
	cfg := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     config.Endpoint.TokenURL,
	}

	token, err := cfg.Token(context.Background())
	if err != nil {
		return validate.NewRequestError(err, http.StatusInternalServerError)
	}

	return web.Respond(ctx, w, token, http.StatusOK)
}
