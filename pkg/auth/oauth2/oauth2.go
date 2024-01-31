package oauth2

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"

	"github.com/lovehotel24/auth-service/pkg/foundation/validate"
	"github.com/lovehotel24/auth-service/pkg/model/user"
	"github.com/lovehotel24/auth-service/pkg/sys/database"
)

//type Server struct {
//	Srv *server.Server
//	mgr *manage.Manager
//	cs  *store.ClientStore
//}

type OAuthServer struct {
	Srv *server.Server
}

func NewOauthServer(userStore user.Store) *OAuthServer {
	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:8080",
	})

	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	srv.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler(userStore))

	return &OAuthServer{srv}
}

//// SetPasswordAuthorizationHandler OAuth Password Grant Handler
//func (as *OAuthServer) SetPasswordAuthorizationHandler(handler server.PasswordAuthorizationHandler) {
//	as.Srv.PasswordAuthorizationHandler = handler
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

//func main() {
//	fmt.Println("I'm NewOauthServer")
//	clientStore := store.NewClientStore()
//	clientStore.Set("222222", &models.Client{
//		ID:     "222222",
//		Secret: "22222222",
//		Domain: "http://localhost:8080",
//	})
//
//	manager := manage.NewDefaultManager()
//
//	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
//
//	// token store
//	manager.MustTokenStorage(store.NewMemoryTokenStore())
//	manager.MapAccessGenerate(generates.NewAccessGenerate())
//	manager.MapClientStorage(clientStore)
//
//	srv := server.NewServer(server.NewConfig(), manager)
//
//	srv.SetPasswordAuthorizationHandler(PasswordAuthorizeHandler)
//	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
//
//	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
//		log.Println("Internal Error:", err.Error())
//		return
//	})
//
//	srv.SetResponseErrorHandler(func(re *errors.Response) {
//		log.Println("Response Error:", re.Error.Error())
//	})
//
//	http.HandleFunc("/v1/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
//
//		store, err := session.Start(r.Context(), w, r)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		var form url.Values
//		if v, ok := store.Get("ReturnUri"); ok {
//			form = v.(url.Values)
//		}
//		r.Form = form
//
//		store.Delete("ReturnUri")
//		store.Save()
//
//		err = srv.HandleAuthorizeRequest(w, r)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//		}
//	})
//
//	http.HandleFunc("/v1/oauth/token", func(w http.ResponseWriter, r *http.Request) {
//
//		err := srv.HandleTokenRequest(w, r)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//		}
//	})
//
//	http.HandleFunc("/v1/oauth/test", func(w http.ResponseWriter, r *http.Request) {
//
//		token, err := srv.ValidationBearerToken(r)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		data := map[string]interface{}{
//			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
//			"client_id":  token.GetClientID(),
//			"user_id":    token.GetUserID(),
//		}
//		e := json.NewEncoder(w)
//		e.SetIndent("", "  ")
//		e.Encode(data)
//	})
//
//	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8081), nil))
//
//	//fmt.Println("I'm client Store")
//	//s := &Server{
//	//	mgr: manage.NewDefaultManager(),
//	//	cs:  clientStore,
//	//}
//
//	//s.mgr.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
//	//
//	//// token store
//	//s.mgr.MustTokenStorage(store.NewMemoryTokenStore())
//
//	// generate jwt access token
//	// s.mgr.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
//	//s.mgr.MapAccessGenerate(generates.NewAccessGenerate())
//	//
//	//s.mgr.MapClientStorage(s.cs)
//	//
//	//s.Srv = server.NewServer(server.NewConfig(), s.mgr)
//	//
//	//fmt.Println("I'm NewServer")
//	//s.Srv.SetPasswordAuthorizationHandler(PasswordAuthorizeHandler)
//	//
//	//fmt.Println("I'm PassAuth")
//	//s.Srv.SetUserAuthorizationHandler(userAuthorizeHandler)
//	//fmt.Println("I'm userAuth")
//	//s.Srv.SetAllowGetAccessRequest(true)
//	//s.Srv.SetClientInfoHandler(server.ClientFormHandler)
//	//fmt.Println("I'm clientInfo")
//	//
//	//s.Srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
//	//	log.Println("Internal Error:", err.Error())
//	//	return
//	//})
//	//
//	//s.Srv.SetResponseErrorHandler(func(re *errors.Response) {
//	//	log.Println("Response Error:", re.Error.Error())
//	//})
//	//fmt.Println("I'm without error")
//	//
//	//return s
//}

//func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
//	fmt.Println("I'm userAuthorizeHandler", r.RequestURI)
//	store, err := session.Start(r.Context(), w, r)
//	if err != nil {
//		return
//	}
//
//	uid, ok := store.Get("LoggedInUserID")
//	if !ok {
//
//		store.Set("ReturnUri", r.Form)
//		store.Save()
//
//		//w.Header().Set("Location", "/v1/oauth/pwd")
//		w.WriteHeader(http.StatusFound)
//		return
//	}
//
//	userID = uid.(string)
//	store.Delete("LoggedInUserID")
//	store.Save()
//	return
//}

//func PasswordAuthorizeHandler(ctx context.Context, clientID, phone, password string) (userID string, err error) {
//	fmt.Println("I'm PasswordAuthorizeHandler ", "- I'm clientID ", clientID, "- I'm phone ", phone)
//	db, err := database.Open(database.Config{
//		User:         "postgres",
//		Password:     "postgres",
//		Host:         "localhost",
//		Name:         "users",
//		MaxIdleConns: 0,
//		MaxOpenConns: 0,
//		DisableTLS:   true,
//	})
//	log, err := logger.New("AUTH")
//
//	userAuth := &userAuth{
//		User: user.NewStore(log, db),
//	}
//
//	//v, err := web.GetValues(ctx)
//	//if err != nil {
//	//	return "", web.NewShutdownError("web value missing from context")
//	//}
//	now := time.Now()
//
//	userID, err = userAuth.User.Authenticate(ctx, now, phone, password)
//	fmt.Println("I'm userID ", userID)
//	if err != nil {
//		switch validate.Cause(err) {
//		case database.ErrNotFound:
//			return "", validate.NewRequestError(err, http.StatusNotFound)
//		case database.ErrAuthenticationFailure:
//			return "", validate.NewRequestError(err, http.StatusUnauthorized)
//		default:
//			return "", fmt.Errorf("authenticating: %w", err)
//		}
//	}
//
//	return userID, nil
//}
