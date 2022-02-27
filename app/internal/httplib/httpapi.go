package httplib

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

//go:generate easyjson

const (
	IDRequest = "id"
)

type (
	//easyjson:json
	HttpError struct {
		Error string `json:"error"`
	}
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	setContentTypeApplicationJson(w)

	w.WriteHeader(statusCode)

	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			logrus.Errorf("write http response: %s", err.Error())
		}
	}
}

func WriteError(w http.ResponseWriter, url, err string, statusCode int) {
	if statusCode >= 500 && statusCode < 600 {
		logrus.Errorf("http error: %s, %s", url, err)
	} else if statusCode == http.StatusBadRequest {
		logrus.Infof("http info: %s, %s", url, err)
	} else {
		logrus.Warnf("http warning: %s, %s", url, err)
	}

	setContentTypeApplicationJson(w)

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(&HttpError{Error: err}); err != nil {
		logrus.Errorf("error json encode: %s", err.Error())
	}
}

func UnmarshalAndValidate(out interface{}, r io.Reader, validate *validator.Validate) (err error) {
	if err = json.NewDecoder(r).Decode(&out); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	if err = validate.Struct(out); err != nil {
		return fmt.Errorf("validate request: %w", err)
	}

	return nil
}

func IDFromRequest(r *http.Request) (userID int64, err error) {
	id := mux.Vars(r)[IDRequest]
	if userID, err = strconv.ParseInt(id, 10, 64); err != nil {
		return 0, fmt.Errorf("parse id: %w", err)
	}

	return userID, nil
}

func setContentTypeApplicationJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
