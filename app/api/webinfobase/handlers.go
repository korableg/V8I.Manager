package webinfobase

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/korableg/V8I.Manager/app/internal/httplib"
	"io/ioutil"
	"net/http"
)

type (
	Handlers struct {
		validate *validator.Validate
		service  Service
	}
)

func NewHandlers(service Service, validate *validator.Validate) (*Handlers, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}

	if validate == nil {
		return nil, errors.New("validator is nil")
	}

	h := &Handlers{validate: validate}

	return h, nil
}

func (h *Handlers) Register(r *mux.Router) *mux.Router {
	infoBasesRouter := r.PathPrefix("/WebCommonInfoBases").Subrouter()

	infoBasesRouter.HandleFunc("", h.Head).Methods("HEAD")
	infoBasesRouter.HandleFunc("", h.GetWSDL).Queries("wsdl", "").Methods("GET")
	infoBasesRouter.HandleFunc("", h.WebCommonInfoBasesPost).Methods("POST")
	infoBasesRouter.HandleFunc("/ws.cws", h.WebCommonInfoBasesPost).Methods("POST")

	return r
}

func (h *Handlers) Head(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handlers) GetWSDL(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(h.service.WSDL())
}

func (h *Handlers) WebCommonInfoBasesPost(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = data

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
        <m:CheckInfoBasesResponse xmlns:m="https://titovcode.com/WebCommonInfoBases">
            <m:return xmlns:xs="http://www.w3.org/2001/XMLSchema"
					xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"/>
            <m:Changed xmlns:xs="http://www.w3.org/2001/XMLSchema"
					xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">true</m:Changed>
            <m:URL xmlns:xs="http://www.w3.org/2001/XMLSchema"
					xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">/WebCommonInfoBases</m:URL>
        </m:CheckInfoBasesResponse>
    </soap:Body>
</soap:Envelope>`))
}
