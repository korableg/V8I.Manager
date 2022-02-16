package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/korableg/V8I.Manager/app/api/user"
	"github.com/korableg/V8I.Manager/app/api/user/service"
	"github.com/korableg/V8I.Manager/app/internal/httperror"
	"github.com/sirupsen/logrus"
)

type (
	Params struct {
		Service   service.Service
		Validator *validator.Validate
	}

	Handlers struct {
		service   service.Service
		validator *validator.Validate
	}
)

func NewHandlers(params Params) (*Handlers, error) {
	if params.Service == nil {
		return nil, errors.New("service is nil")
	}

	if params.Validator == nil {
		return nil, errors.New("validator is nil")
	}

	h := &Handlers{
		service:   params.Service,
		validator: params.Validator,
	}

	return h, nil
}

func (h *Handlers) Register(r *mux.Router) *mux.Router {
	r.HandleFunc("/add", h.Add).Methods("POST")

	return r
}

func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	var req user.AddUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		httperror.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Add(r.Context(), req)
	if err != nil {
		httperror.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := user.AddUserResponse{ID: id}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&resp); err != nil {
		logrus.Errorf("json marshal: %s", err.Error())
	}
}
