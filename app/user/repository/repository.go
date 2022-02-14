package repository

import (
	"context"
	"github.com/korableg/V8I.Manager/app/user"
)

type Repository interface {
	Add(ctx context.Context, u user.User) (int64, error)
	Get(ctx context.Context, ID int64) (user.User, error)
	GetByName(ctx context.Context, name string) (user.User, error)
	GetList(ctx context.Context) ([]user.User, error)
	Update(ctx context.Context, u user.User) error
	Delete(ctx context.Context, ID int64) error
}
