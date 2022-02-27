package user

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/korableg/V8I.Manager/app/internal/httplib"
	"net/http"
)

type (
	Handlers struct {
		service  Service
		validate *validator.Validate
	}
)

func NewHandlers(svc Service, validate *validator.Validate) (*Handlers, error) {
	if svc == nil {
		return nil, errors.New("service is nil")
	}

	if validate == nil {
		return nil, errors.New("validator is nil")
	}

	h := &Handlers{
		service:  svc,
		validate: validate,
	}

	return h, nil
}

func (h *Handlers) Register(r *mux.Router) *mux.Router {
	requestWithUserID := fmt.Sprintf("/{%s:[0-9]+}", httplib.IDRequest)

	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", h.Add).Methods("POST")
	userRouter.HandleFunc(requestWithUserID, h.Update).Methods("PUT")
	userRouter.HandleFunc(requestWithUserID, h.Get).Methods("GET")
	userRouter.HandleFunc("", h.GetList).Methods("GET")
	userRouter.HandleFunc(requestWithUserID, h.Delete).Methods("DELETE")

	return r
}

func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	var (
		req AddUserRequest
		err error
	)

	if err = httplib.UnmarshalAndValidate(&req, r.Body, h.validate); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Add(r.Context(), req)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusCreated, &AddUserResponse{ID: id})
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	userID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Get(r.Context(), userID)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, &user)
}

func (h *Handlers) GetList(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetList(r.Context())
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, &user)
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	var (
		req UpdateUserRequest
		err error
	)

	if err = httplib.UnmarshalAndValidate(&req, r.Body, h.validate); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	req.ID = userID

	if err = h.service.Update(r.Context(), req); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, nil)
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.Delete(r.Context(), userID); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, nil)
}
