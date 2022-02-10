package user

import (
	"context"
)

type Repository interface {
	Add(ctx context.Context, u User) (int64, error)
	GetList(ctx context.Context) ([]User, error)
}
