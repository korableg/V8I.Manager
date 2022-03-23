package webinfobase

import "github.com/google/uuid"

type (
	CheckInfoBaseResponse struct {
		URL     string
		Changed bool
	}

	GetInfoBasesResponse struct {
		ClientID uuid.UUID
		Code     uuid.UUID
		Text     string
	}
)
