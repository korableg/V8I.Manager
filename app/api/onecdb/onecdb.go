package onecdb

//go:generate easyjson

import (
	"context"
	"github.com/google/uuid"
)

const (
	ClientConnectionSpeedNormal = "normal"
	ClientConnectionSpeedLow    = "low"

	AppAuto        = "Auto"
	AppThinClient  = "ThinClient"
	AppThickClient = "ThickClient"
	AppWebClient   = "WebClient"

	WAEnabled  = 1
	WADisabled = 0
)

type (
	Repository interface {
		Add(ctx context.Context, db DB) (int64, error)
		Get(ctx context.Context, id int64) (DB, error)
		GetList(ctx context.Context) ([]DB, error)
		Update(ctx context.Context, db DB) error
		Delete(ctx context.Context, id int64) error
	}

	DBCollector interface {
		Collect(db ...DB) error
	}

	V8IBuilder interface {
		BuildV8I(ctx context.Context) (string, error)
	}

	Service interface {
		Add(ctx context.Context, reqDB AddDBRequest) (int64, error)
		Get(ctx context.Context, ID int64) (DB, error)
		GetList(ctx context.Context) ([]DB, error)
		Update(ctx context.Context, reqDB UpdateDBRequest) error
		Delete(ctx context.Context, ID int64) error
	}

	// DB Description of v8i file: https://its.1c.ru/db/v838doc#bookmark:adm:TI000000368
	//easyjson:json
	DB struct {
		ID                    int64     `json:"id" validate:"required"`
		UUID                  uuid.UUID `json:"uuid" validate:"required"`
		Name                  string    `json:"name" validate:"required"`
		Connect               string    `json:"connect" validate:"required"`
		OrderInList           int64     `json:"order_in_list" validate:"required"`
		OrderInTree           int64     `json:"order_in_tree" validate:"required"`
		Folder                string    `json:"folder" validate:"required"`
		ClientConnectionSpeed string    `json:"client_connection_speed" validate:"oneof='normal' 'low'"`
		App                   string    `json:"app" validate:"oneof='Auto' 'ThinClient' 'ThickClient' 'WebClient'"`
		WA                    int64     `json:"wa" validate:"gte=0,lte=1"` //0 - всегда использовать аутентификацию с помощью логина/пароля, 1 -пытаться выполнить аутентификацию средствами ОС. Если выполнено неудачно, запрашивается логин/пароль,
		Version               string    `json:"version"`
		WebCommonInfoBaseURL  string    `json:"web_common_info_base_url"`
		AdditionalParameters  string    `json:"additional_parameters"`
	}

	//easyjson:json
	AddDBRequest struct {
		UUID                  string `json:"uuid" validate:"required,uuid"`
		Name                  string `json:"name" validate:"required"`
		Connect               string `json:"connect" validate:"required"`
		OrderInList           int64  `json:"order_in_list" validate:"required"`
		OrderInTree           int64  `json:"order_in_tree" validate:"required"`
		Folder                string `json:"folder" validate:"required"`
		ClientConnectionSpeed string `json:"client_connection_speed" validate:"oneof='normal' 'low'"`
		App                   string `json:"app" validate:"oneof='Auto' 'ThinClient' 'ThickClient' 'WebClient'"`
		WA                    int64  `json:"wa" validate:"gte=0,lte=1"` //0 - всегда использовать аутентификацию с помощью логина/пароля, 1 -пытаться выполнить аутентификацию средствами ОС. Если выполнено неудачно, запрашивается логин/пароль,
		Version               string `json:"version"`
		AdditionalParameters  string `json:"additional_parameters"`
	}

	//easyjson:json
	UpdateDBRequest struct {
		ID                    int64  `json:"-"`
		UUID                  string `json:"uuid" validate:"required,uuid4"`
		Name                  string `json:"name" validate:"required"`
		Connect               string `json:"connect" validate:"required"`
		OrderInList           int64  `json:"order_in_list" validate:"required"`
		OrderInTree           int64  `json:"order_in_tree" validate:"required"`
		Folder                string `json:"folder" validate:"required"`
		ClientConnectionSpeed string `json:"client_connection_speed" validate:"oneof='normal' 'low'"`
		App                   string `json:"app" validate:"oneof='Auto' 'ThinClient' 'ThickClient' 'WebClient'"`
		WA                    int64  `json:"wa" validate:"gte=0,lte=1"` //0 - всегда использовать аутентификацию с помощью логина/пароля, 1 -пытаться выполнить аутентификацию средствами ОС. Если выполнено неудачно, запрашивается логин/пароль,
		Version               string `json:"version"`
		AdditionalParameters  string `json:"additional_parameters"`
	}

	//easyjson:json
	AddDBResponse struct {
		ID int64 `json:"id"`
	}
)
