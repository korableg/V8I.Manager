package httperror

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

//go:generate easyjson

type (
	//easyjson:json
	HttpError struct {
		Error string `json:"error"`
	}
)

func WriteError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	if statusCode >= 500 && statusCode < 600 {
		logrus.Errorf("http error: %s, %s", r.RequestURI, err.Error())
	} else if statusCode == http.StatusBadRequest {
		logrus.Infof("http info: %s, %s", r.RequestURI, err.Error())
	} else {
		logrus.Warnf("http warning: %s, %s", r.RequestURI, err.Error())
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(&HttpError{Error: err.Error()}); err != nil {
		logrus.Errorf("error json encode: %s", err.Error())
	}
}
