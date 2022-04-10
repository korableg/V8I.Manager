package onecserver

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

func (h *Handlers) Register(r *mux.Router) {
	requestWithServerID := fmt.Sprintf("/{%s:[0-9]+}", httplib.IDRequest)

	serverRouter := r.PathPrefix("/servers").Subrouter()
	serverRouter.HandleFunc("", h.Add).Methods("POST")
	serverRouter.HandleFunc(requestWithServerID, h.Get).Methods("GET")
	serverRouter.HandleFunc("", h.GetList).Methods("GET")
	serverRouter.HandleFunc(requestWithServerID, h.Update).Methods("PUT")
	serverRouter.HandleFunc(fmt.Sprintf("%s/switch-watching", requestWithServerID), h.SwitchWatching).Methods("POST")
	serverRouter.HandleFunc(requestWithServerID, h.Delete).Methods("DELETE")

	return r
}

func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	var (
		req AddServerRequest
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

	httplib.WriteJSONResponse(w, http.StatusCreated, &AddServerResponse{ID: id})
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	serverID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	server, err := h.service.Get(r.Context(), serverID)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, &server)
}

func (h *Handlers) GetList(w http.ResponseWriter, r *http.Request) {
	servers, err := h.service.GetList(r.Context())
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, &servers)
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	var (
		req UpdateServerRequest
		err error
	)

	if err = httplib.UnmarshalAndValidate(&req, r.Body, h.validate); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	serverID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	req.ID = serverID

	if err = h.service.Update(r.Context(), req); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, nil)
}

func (h *Handlers) SwitchWatching(w http.ResponseWriter, r *http.Request) {
	serverID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	watching, err := h.service.SwitchWatching(r.Context(), serverID)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, &SwitchWatchingResponse{Watching: watching})
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	serverID, err := httplib.IDFromRequest(r)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.Delete(r.Context(), serverID); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	httplib.WriteJSONResponse(w, http.StatusOK, nil)
}
