package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-session/session"

	"github.com/lovehotel24/auth-service/pkg/foundation/validate"
	"github.com/lovehotel24/auth-service/pkg/foundation/web"
	"github.com/lovehotel24/auth-service/pkg/model/user"
	"github.com/lovehotel24/auth-service/pkg/sys/database"
)

type Handlers struct {
	User user.Store
}

// Query returns a list of users with paging.
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page := web.Param(r, "page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid page format [%s]", page), http.StatusBadRequest)
	}

	rows := web.Param(r, "rows")
	rowsPerPage, err := strconv.Atoi(rows)
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid rows format [%s]", rows), http.StatusBadRequest)
	}

	users, err := h.User.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for users: %w", err)
	}

	return web.Respond(ctx, w, users, http.StatusOK)
}

// QueryByID returns a user by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//claims, err := auth.GetClaims(ctx)
	//if err != nil {
	//	return errors.New("claims missing from context")
	//}

	id := web.Param(r, "id")
	usr, err := h.User.QueryByID(ctx, id)
	if err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidID:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

// Create adds a new user to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var nu user.NewUser
	//b, err := io.ReadAll(r.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(string(b))

	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	usr, err := h.User.Create(ctx, nu, v.Now)
	if err != nil {
		return fmt.Errorf("user[%+v]: %w", &usr, err)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

// Update updates a user in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	//claims, err := auth.GetClaims(ctx)
	//if err != nil {
	//	return errors.New("claims missing from context")
	//}

	var upd user.UpdateUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	id := web.Param(r, "id")
	if err := h.User.Update(ctx, id, upd, v.Now); err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidID:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("ID[%s] User[%+v]: %w", id, &upd, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a user from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//claims, err := auth.GetClaims(ctx)
	//if err != nil {
	//	return errors.New("claims missing from context")
	//}

	id := web.Param(r, "id")
	if err := h.User.Delete(ctx, id); err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidID:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h Handlers) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	store, err := session.Start(r.Context(), w, r)
	if err != nil {

		return validate.NewRequestError(err, http.StatusInternalServerError)
	}

	phone, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide phone and password in Basic auth")
		return validate.NewRequestError(err, http.StatusUnauthorized)
	}

	userID, err := h.User.Authenticate(ctx, v.Now, phone, pass)
	if err != nil {
		switch validate.Cause(err) {
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrAuthenticationFailure:
			return validate.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("authenticating: %w", err)
		}
	}
	store.Set("LoggedInUserID", userID)
	store.Save()

	//var tkn struct {
	//	Token string `json:"token"`
	//}
	//tkn.Token, err = h.Auth.GenerateToken(claims)
	//if err != nil {
	//	return fmt.Errorf("generating token: %w", err)
	//}

	w.Header().Set("Location", "/v1/oauth/authorize")

	status := struct {
		Status string
	}{
		Status: "Login Success",
	}

	return web.Respond(ctx, w, status, http.StatusFound)
}

func (h Handlers) PasswordAuthorizeHandler(ctx context.Context, clientID, phone, password string) (userID string, err error) {

	now := time.Now()

	userID, err = h.User.Authenticate(ctx, now, phone, password)
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
