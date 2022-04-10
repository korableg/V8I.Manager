package onecserver

import (
	"context"
	"errors"
	"io"
)

//go:generate easyjson

type (
	Repository interface {
		Add(ctx context.Context, u Server) (int64, error)
		Get(ctx context.Context, ID int64) (Server, error)
		GetList(ctx context.Context) ([]Server, error)
		Update(ctx context.Context, u Server) error
		UpdateHash(ctx context.Context, ID int64, hash string) error
		UpdateWatch(ctx context.Context, ID int64, watch bool) error
		Delete(ctx context.Context, ID int64) error
	}

	Service interface {
		Add(ctx context.Context, u AddServerRequest) (int64, error)
		Get(ctx context.Context, ID int64) (Server, error)
		GetList(ctx context.Context) ([]Server, error)
		Update(ctx context.Context, u UpdateServerRequest) error
		SwitchWatching(ctx context.Context, ID int64) (bool, error)
		Delete(ctx context.Context, ID int64) error
		io.Closer
	}

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

	//easyjson:json
	AddServerResponse struct {
		ID int64 `json:"id"`
	}

	//easyjson:json
	SwitchWatchingResponse struct {
		Watching bool `json:"watching"`
	}
)

var (
	ErrServerNotFound = errors.New("server not found")
)
