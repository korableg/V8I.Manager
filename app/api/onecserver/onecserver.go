package onecserver

import "errors"

//go:generate easyjson

type (
	//easyjson:json
	Server struct {
		ID      int64  `json:"id"`
		Name    string `json:"name"`
		LSTPath string `json:"lst_path"`
		Watch   bool   `json:"watch"`
	}

	//easyjson:json
	AddServerRequest struct {
		Name    string `json:"name"`
		LSTPath string `json:"lst_path" validate:"required,file"`
	}

	//easyjson:json
	UpdateServerRequest struct {
		ID      int64
		Name    string `json:"name" validate:"required"`
		LSTPath string `json:"lst_path" validate:"required,file"`
	}
)

var (
	ErrServerNotFound = errors.New("server not found")
)
