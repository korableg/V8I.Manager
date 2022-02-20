package onecdb

import (
	"context"
	"github.com/google/uuid"
)

type (
	Repository interface {
		Add(ctx context.Context, db DB) error
		Get(ctx context.Context, id uuid.UUID) (DB, error)
		GetList(ctx context.Context) ([]DB, error)
		Update(ctx context.Context, db DB) error
		Delete(ctx context.Context, id uuid.UUID) error
	}
)
