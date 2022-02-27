package engine

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/korableg/V8I.Manager/app/api/onecdb"
	"github.com/korableg/V8I.Manager/app/api/user"
	"github.com/korableg/V8I.Manager/app/api/user/auth"
	"github.com/korableg/V8I.Manager/app/internal/config"
	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
	"github.com/korableg/V8I.Manager/app/internal/transport/httpserver"
)

type (
	Engine struct {
		httpSrvr *httpserver.HttpServer
		sqliteDB *sqlitedb.SqliteDB
	}
)

func NewEngine(cfgPath string) (*Engine, error) {
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	validate := validator.New()

	sdb, err := sqlitedb.NewSqliteDB(cfg.Sqlite)
	if err != nil {
		return nil, fmt.Errorf("init sqlite db: %w", err)
	}

	userRepo, err := user.NewSqliteRepository(sdb)
	if err != nil {
		return nil, fmt.Errorf("init user repository: %w", err)
	}

	userHds, err := initUserHandlers(userRepo, validate)
	if err != nil {
		return nil, fmt.Errorf("user handlers: %w", err)
	}

	authHds, err := initAuthHandlers(userRepo, cfg.Auth, validate)
	if err != nil {
		return nil, fmt.Errorf("auth handlers: %w", err)
	}

	dbHds, err := initDBHandlers(sdb, validate)
	if err != nil {
		return nil, fmt.Errorf("onecdb handlers: %w", err)
	}

	httpSrvr := httpserver.NewHttpServer(
		cfg.Http,
		httpserver.WithApiMiddleware(authHds.Middleware()),
		httpserver.WithApiHandlers(userHds),
		httpserver.WithApiHandlers(dbHds),
		httpserver.WithHandlers(authHds),
	)

	return &Engine{
		httpSrvr: httpSrvr,
		sqliteDB: sdb,
	}, nil
}

func (en *Engine) Start() error {
	if err := en.httpSrvr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve http: %w", err)
	}

	return nil
}

func (en *Engine) Shutdown(ctx context.Context) error {
	if err := en.httpSrvr.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	if err := en.sqliteDB.Close(); err != nil {
		return fmt.Errorf("close db: %w", err)
	}

	return nil
}

func initUserHandlers(userRepo user.Repository, validate *validator.Validate) (*user.Handlers, error) {
	userSvc, err := user.NewService(userRepo)
	if err != nil {
		return nil, fmt.Errorf("init user service: %w", err)
	}

	if _, err = userSvc.Add(context.Background(), user.AddUserRequest{
		Name:     "admin",
		Password: "admin",
		Role:     "admin",
	}); err != nil && !errors.Is(err, user.ErrUserAlreadyCreated) {
		return nil, fmt.Errorf("init admin user: %w", err)
	}

	userHds, err := user.NewHandlers(userSvc, validate)
	if err != nil {
		return nil, fmt.Errorf("init user handlers: %w", err)
	}

	return userHds, nil
}

func initAuthHandlers(userRepo user.Repository, authCfg auth.Config, validate *validator.Validate) (*auth.Handlers, error) {
	authSvc := auth.NewAuth(userRepo, authCfg)

	authHds, err := auth.NewHandlers(authSvc, validate)
	if err != nil {
		return nil, fmt.Errorf("init auth handlers: %w", err)
	}

	return authHds, nil
}

func initDBHandlers(sdb *sqlitedb.SqliteDB, validate *validator.Validate) (*onecdb.Handlers, error) {
	dbRepo, err := onecdb.NewSqliteRepository(sdb)
	if err != nil {
		return nil, fmt.Errorf("init onecdb repository: %w", err)
	}

	dbService, err := onecdb.NewService(dbRepo)
	if err != nil {
		return nil, fmt.Errorf("init onecdb service: %w", err)
	}

	dbHandlers, err := onecdb.NewHandlers(dbService, validate)
	if err != nil {
		return nil, fmt.Errorf("init onecdb handlers: %w", err)
	}

	return dbHandlers, nil
}
