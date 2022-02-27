package onecdb

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/korableg/V8I.Manager/app/internal/httplib"
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
	requestWithDBID := fmt.Sprintf("/{%s:[0-9]+}", httplib.IDRequest)

	dbRouter := r.PathPrefix("/dbs").Subrouter()
	dbRouter.HandleFunc("", h.Add).Methods("POST")
	dbRouter.HandleFunc(requestWithDBID, h.Update).Methods("PUT")
	dbRouter.HandleFunc(requestWithDBID, h.Get).Methods("GET")
	dbRouter.HandleFunc("", h.GetList).Methods("GET")
	dbRouter.HandleFunc(requestWithDBID, h.Delete).Methods("DELETE")

	return r
}

func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	var (
		req AddDBRequest
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

	httplib.WriteJSONResponse(w, http.StatusCreated, &AddDBResponse{ID: id})
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	dbID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := h.service.Get(r.Context(), dbID)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, &db)
}

func (h *Handlers) GetList(w http.ResponseWriter, r *http.Request) {
	dbs, err := h.service.GetList(r.Context())
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, &dbs)
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	var (
		req UpdateDBRequest
		err error
	)

	if err = httplib.UnmarshalAndValidate(&req, r.Body, h.validate); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	dbID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	req.ID = dbID

	if err = h.service.Update(r.Context(), req); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, nil)
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	dbID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.Delete(r.Context(), dbID); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, nil)
}
