package repository

import (
	"context"
	"github.com/korableg/V8I.Manager/app/user"
)

type Repository interface {
	Add(ctx context.Context, u user.User) (int64, error)
	GetList(ctx context.Context) ([]user.User, error)
}
