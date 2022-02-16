package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/korableg/V8I.Manager/app/api/user"
	"github.com/korableg/V8I.Manager/app/api/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		Add(ctx context.Context, u user.AddUserRequest) (int64, error)
		Get(ctx context.Context, ID int64) (user.User, error)
		GetList(ctx context.Context) ([]user.User, error)
		Update(ctx context.Context, u user.UpdateUserRequest) error
		Delete(ctx context.Context, ID int64) error
	}

	service struct {
		userRepo repository.Repository
	}
)

func NewService(userRepo repository.Repository) (*service, error) {
	if userRepo == nil {
		return nil, errors.New("user repository is not defined")
	}

	s := &service{
		userRepo: userRepo,
	}

	return s, nil
}

func (s *service) Add(ctx context.Context, uReq user.AddUserRequest) (int64, error) {
	passwordHash, err := hashPassword(uReq.Password)
	if err != nil {
		return 0, fmt.Errorf("hashPassword: %w", err)
	}

	u := user.User{
		Name:         uReq.Name,
		PasswordHash: passwordHash,
		Token:        uReq.Token,
		Role:         uReq.Role,
	}

	id, err := s.userRepo.Add(ctx, u)
	if err != nil {
		return 0, fmt.Errorf("add to store: %w", err)
	}

	return id, nil
}

func (s *service) Get(ctx context.Context, ID int64) (user.User, error) {
	u, err := s.Get(ctx, ID)
	if err != nil {
		return user.User{}, fmt.Errorf("get from store: %w", err)
	}

	return u, nil
}

func (s *service) GetList(ctx context.Context) ([]user.User, error) {
	users, err := s.userRepo.GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get from store: %w", err)
	}

	return users, nil
}

func (s *service) Update(ctx context.Context, uReq user.UpdateUserRequest) error {
	u := user.User{
		ID:    uReq.ID,
		Name:  uReq.Name,
		Token: uReq.Token,
		Role:  uReq.Role,
	}

	if err := s.userRepo.Update(ctx, u); err != nil {
		return fmt.Errorf("update user in store: %w", err)
	}

	return nil
}

func (s service) Delete(ctx context.Context, ID int64) error {
	if err := s.userRepo.Delete(ctx, ID); err != nil {
		return fmt.Errorf("delete user from store: %w", err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}

	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
