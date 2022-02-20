package auth

//go:generate easyjson

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/korableg/V8I.Manager/app/api/user"
)

type (
	//easyjson:json
	UserClaims struct {
		jwt.StandardClaims

		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	Config struct {
		Secret string `yaml:"jwt_secret""`
	}

	Auth interface {
		SignIn(ctx context.Context, req user.SignInRequest) (string, time.Time, error)
		GetUserFromToken(ctx context.Context, token string) (user.User, error)
	}

	auth struct {
		userRepo  user.Repository
		jwtSecret []byte
	}
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidPassword = errors.New("invalid password")
)

func NewAuth(userRepo user.Repository, jwtCfg Config) *auth {
	a := &auth{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtCfg.Secret),
	}

	return a
}

func (a *auth) SignIn(ctx context.Context, req user.SignInRequest) (string, time.Time, error) {
	u, err := a.userRepo.GetByName(ctx, req.Name)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return "", time.Time{}, err
		}

		return "", time.Time{}, fmt.Errorf("get user by name: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", time.Time{}, ErrInvalidPassword
		}

		return "", time.Time{}, fmt.Errorf("compare password: %w", err)
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(7 * 86400 * time.Second)

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
		ID:   u.ID,
		Name: u.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	j, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("jwt signed string: %w", err)
	}

	return j, expiresAt, nil
}

func (a *auth) GetUserFromToken(ctx context.Context, jwtToken string) (user.User, error) {
	claims := UserClaims{}

	token, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return user.User{}, ErrInvalidToken
		}

		return user.User{}, fmt.Errorf("token error: %w", err)
	}

	if !token.Valid {
		return user.User{}, ErrInvalidToken
	}

	u, err := a.userRepo.Get(ctx, claims.ID)
	if err != nil {
		return user.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return u, nil
}
