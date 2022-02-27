package onecdb

import (
	"context"
)

type (
	Repository interface {
		Add(ctx context.Context, db DB) (int64, error)
		Get(ctx context.Context, id int64) (DB, error)
		GetList(ctx context.Context) ([]DB, error)
		Update(ctx context.Context, db DB) error
		Delete(ctx context.Context, id int64) error
	}
)
