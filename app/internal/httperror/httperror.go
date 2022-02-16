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

func WriteError(w http.ResponseWriter, url, err string, statusCode int) {
	if statusCode >= 500 && statusCode < 600 {
		logrus.Errorf("http error: %s, %s", url, err)
	} else if statusCode == http.StatusBadRequest {
		logrus.Infof("http info: %s, %s", url, err)
	} else {
		logrus.Warnf("http warning: %s, %s", url, err)
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(&HttpError{Error: err}); err != nil {
		logrus.Errorf("error json encode: %s", err.Error())
	}
}
