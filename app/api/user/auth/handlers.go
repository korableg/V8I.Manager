package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/korableg/V8I.Manager/app/api/user"
	"github.com/korableg/V8I.Manager/app/internal/httplib"
)

const (
	jwtTokenCookieName = "jwt_auth"
)

type (
	Handlers struct {
		authSvc  Auth
		validate *validator.Validate
	}
)

func NewHandlers(svc Auth, validate *validator.Validate) (*Handlers, error) {
	if svc == nil {
		return nil, errors.New("auth service is nil")
	}

	if validate == nil {
		return nil, errors.New("validator is nil")
	}

	h := &Handlers{
		authSvc:  svc,
		validate: validate,
	}

	return h, nil
}

func (h *Handlers) Register(r *mux.Router) *mux.Router {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signin", h.SignIn).Methods("POST")
	authRouter.HandleFunc("/signout", h.SignOut).Methods("GET")

	return r
}

func (h *Handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		req user.SignInRequest
		err error
	)

	if err = httplib.UnmarshalAndValidate(&req, r.Body, h.validate); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	jwtToken, jwtExpires, err := h.authSvc.SignIn(r.Context(), req)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, ErrInvalidPassword) {
			httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
			return
		}

		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     jwtTokenCookieName,
		Value:    jwtToken,
		Path:     "/",
		Expires:  jwtExpires,
		HttpOnly: true,
	})

	httplib.WriteJSONResponse(w, http.StatusOK, nil)
}

func (h *Handlers) SignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     jwtTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	httplib.WriteJSONResponse(w, http.StatusOK, nil)
}

func (h *Handlers) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtCookie, err := r.Cookie(jwtTokenCookieName)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					httplib.WriteError(w, r.RequestURI, fmt.Sprintf("unathorized: %s", err.Error()), http.StatusUnauthorized)
					return
				}

				httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
				return
			}

			u, err := h.authSvc.GetUserFromToken(r.Context(), jwtCookie.Value)
			if err != nil {
				if errors.Is(err, ErrInvalidToken) || errors.Is(err, user.ErrUserNotFound) {
					httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusUnauthorized)
					return
				}

				httplib.WriteError(w, r.RequestURI, fmt.Sprintf("get user from token: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), user.CtxKey, u)))
		})
	}
}
