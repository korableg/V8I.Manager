package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/korableg/V8I.Manager/app/internal/httperror"
	"github.com/korableg/V8I.Manager/app/user"
	"github.com/korableg/V8I.Manager/app/user/service"
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

func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	var req user.AddUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteError(w, r, fmt.Errorf("parse json: %w", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		httperror.WriteError(w, r, fmt.Errorf("validate request: %w", err), http.StatusBadRequest)
		return
	}

	id, err := h.service.Add(r.Context(), req)
	if err != nil {
		httperror.WriteError(w, r, fmt.Errorf("add user: %w", err), http.StatusInternalServerError)
		return
	}

	resp := user.AddUserResponse{ID: id}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&resp); err != nil {
		logrus.Errorf("json marshal: %s", err.Error())
	}
}
