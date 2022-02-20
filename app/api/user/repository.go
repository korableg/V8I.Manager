package user

import (
	"context"
	"errors"
)

type Repository interface {
	Add(ctx context.Context, u User) (int64, error)
	Get(ctx context.Context, ID int64) (User, error)
	GetByName(ctx context.Context, name string) (User, error)
	GetList(ctx context.Context) ([]User, error)
	Update(ctx context.Context, u User) error
	Delete(ctx context.Context, ID int64) error
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyCreated = errors.New("user is already created")
)
