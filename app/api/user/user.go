package user

import (
	"context"
	"errors"
)

//go:generate easyjson

type (
	ctxKey struct{}

	Service interface {
		Add(ctx context.Context, u AddUserRequest) (int64, error)
		Get(ctx context.Context, ID int64) (User, error)
		GetList(ctx context.Context) ([]User, error)
		Update(ctx context.Context, u UpdateUserRequest) error
		Delete(ctx context.Context, ID int64) error
	}

	Repository interface {
		Add(ctx context.Context, u User) (int64, error)
		Get(ctx context.Context, ID int64) (User, error)
		GetByName(ctx context.Context, name string) (User, error)
		GetList(ctx context.Context) ([]User, error)
		Update(ctx context.Context, u User) error
		Delete(ctx context.Context, ID int64) error
	}

	//easyjson:json
	AddUserRequest struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required,oneof='admin' 'reader'"`
		Token    string `json:"token"`
	}

	//easyjson:json
	UpdateUserRequest struct {
		ID    int64  `json:"-"`
		Name  string `json:"name" validate:"required"`
		Role  string `json:"role" validate:"required,oneof='admin' 'reader'"`
		Token string `json:"token"`
	}

	//easyjson:json
	AddUserResponse struct {
		ID int64 `json:"id"`
	}

	//easyjson:json
	SignInRequest struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	//easyjson:json
	User struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		PasswordHash string `json:"-"`
		Token        string `json:"token"`
		Role         string `json:"role"`
	}
)

var (
	CtxKey ctxKey

	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyCreated = errors.New("user is already created")
)
