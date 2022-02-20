package onecdb

import "github.com/google/uuid"

type (
	// Description of v8i file: https://its.1c.ru/db/v838doc#bookmark:adm:TI000000368
	DB struct {
		ID                    uuid.UUID `json:"id" validate:"require"`
		Name                  string    `json:"name" validate:"require"`
		Connect               string
		OrderInList           int64
		OrderInTree           int64
		Folder                string
		ClientConnectionSpeed string `validate:"oneof:'normal' 'low'"`
		App                   string `validate:"oneof:'Auto' 'ThinClient' 'ThickClient' 'WebClient'"`
		WA                    int64  //0 - всегда использовать аутентификацию с помощью логина/пароля, 1 -пытаться выполнить аутентификацию средствами ОС. Если выполнено неудачно, запрашивается логин/пароль,
		Version               string
		WebCommonInfoBaseURL  string
		AdditionalParameters  string
	}
)
