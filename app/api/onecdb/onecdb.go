package onecdb

//go:generate easyjson

import "github.com/google/uuid"

type (
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
