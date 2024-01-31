package oauth

import (
	"context"
	"net/http"
	"time"

	"github.com/lovehotel24/auth-service/pkg/auth/oauth2"
	"github.com/lovehotel24/auth-service/pkg/foundation/validate"
	"github.com/lovehotel24/auth-service/pkg/foundation/web"
)

type Handlers struct {
	Oauth2 *oauth2.OAuthServer
}

func (h Handlers) Authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return h.Oauth2.Srv.HandleAuthorizeRequest(w, r)
}

func (h Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return h.Oauth2.Srv.HandleTokenRequest(w, r)
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	token, err := h.Oauth2.Srv.ValidationBearerToken(r)
	if err != nil {
		return validate.NewRequestError(err, http.StatusBadRequest)
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}

	return web.Respond(ctx, w, data, http.StatusOK)
}

//type Handlers struct {
//	Server *oauth2.Server
//}

//func (h Handlers) Authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//
//	fmt.Println("I'm Authorize /oauth/authorize")
//
//	fmt.Println(r.Context().Value("ReturnUri"))
//
//	//var form url.Values
//	//if v, ok := store.Get("ReturnUri"); ok {
//	//	fmt.Println("I'm ok from url.Values")
//	//	form = v.(url.Values)
//	//	fmt.Println("I'm ok from url.Values", v.(url.Values))
//	//}
//	//r.Form = form
//	//
//	//store.Delete("ReturnUri")
//	//store.Save()
//	return h.Server.Srv.HandleAuthorizeRequest(w, r)
//}
//
//func (h Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//
//	fmt.Println("I'm Token /oauth/token")
//	return h.Server.Srv.HandleTokenRequest(w, r)
//}
//
//func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//
//	token, err := h.Server.Srv.ValidationBearerToken(r)
//	if err != nil {
//		return validate.NewRequestError(err, http.StatusBadRequest)
//	}
//
//	data := map[string]interface{}{
//		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
//		"client_id":  token.GetClientID(),
//		"user_id":    token.GetUserID(),
//	}
//
//	return web.Respond(ctx, w, data, http.StatusOK)
//}
