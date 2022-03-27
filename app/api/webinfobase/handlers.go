package webinfobase

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/korableg/V8I.Manager/app/internal/httplib"
	"github.com/sirupsen/logrus"
	"net/http"
	"text/template"
)

type (
	Handlers struct {
		validate   *validator.Validate
		service    Service
		checkTmplt *template.Template
		getTmplt   *template.Template
	}
)

//go:embed CheckInfoBasesResponse.xml
var checkInfoBasesTmplt string

//go:embed GetInfoBasesResponse.xml
var getInfoBasesTmplt string

func NewHandlers(svc Service, validate *validator.Validate) (*Handlers, error) {
	if validate == nil {
		return nil, errors.New("validator is nil")
	}

	if svc == nil {
		return nil, errors.New("service is nil")
	}

	checkTmplt, err := template.New("checkInfoBasesResponse").Parse(checkInfoBasesTmplt)
	if err != nil {
		return nil, fmt.Errorf("error parsing checkinfobases template: %w", err)
	}

	getTmplt, err := template.New("getInfoBasesResponse").Parse(getInfoBasesTmplt)
	if err != nil {
		return nil, fmt.Errorf("error parsing getinfobases template: %w", err)
	}

	h := &Handlers{
		validate:   validate,
		service:    svc,
		checkTmplt: checkTmplt,
		getTmplt:   getTmplt,
	}

	return h, nil
}

func (h *Handlers) Register(r *mux.Router) *mux.Router {
	infoBasesRouter := r.PathPrefix("/WebCommonInfoBases").Subrouter()

	infoBasesRouter.HandleFunc("", h.GetWSDL).Queries("wsdl", "").Methods("GET")
	infoBasesRouter.HandleFunc("/ws.cws", h.CheckInfoBases).Methods("POST").HeadersRegexp("Soapaction", "CheckInfoBases")
	infoBasesRouter.HandleFunc("/ws.cws", h.GetInfoBases).Methods("POST").HeadersRegexp("Soapaction", "GetInfoBases")

	return r
}

func (h *Handlers) GetWSDL(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(h.service.WSDL()); err != nil {
		logrus.Errorf("failed to write http body: %s", err.Error())
	}
}

func (h *Handlers) CheckInfoBases(w http.ResponseWriter, r *http.Request) {
	req := CheckInfoBasesRequestWrapper{}
	if err := httplib.UnmarshalXMLAndValidate(&req, r.Body, h.validate); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.CheckInfoBases(r.Context(), req.Body.CheckInfoBasesRequest.ID)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = h.checkTmplt.Execute(w, &resp); err != nil {
		logrus.Errorf("failed to write http body: %s", err)
	}
}

func (h *Handlers) GetInfoBases(w http.ResponseWriter, r *http.Request) {
	req := GetInfoBasesRequestWrapper{}
	if err := httplib.UnmarshalXMLAndValidate(&req, r.Body, h.validate); err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetInfoBases(r.Context(), req.Body.GetInfoBasesRequest.ID)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = h.getTmplt.Execute(w, &resp); err != nil {
		logrus.Errorf("failed to write http body: %s", err)
	}
}
