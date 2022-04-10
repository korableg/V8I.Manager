package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/korableg/V8I.Manager/app/api/client"
	"github.com/korableg/V8I.Manager/app/api/onecdb"
	"github.com/korableg/V8I.Manager/app/api/onecserver"
	"github.com/korableg/V8I.Manager/app/api/onecserver/dbbuilder"
	"github.com/korableg/V8I.Manager/app/api/onecserver/watcher"
	"github.com/korableg/V8I.Manager/app/api/user"
	"github.com/korableg/V8I.Manager/app/api/user/auth"
	"github.com/korableg/V8I.Manager/app/api/webinfobase"
	"github.com/korableg/V8I.Manager/app/internal/config"
	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
	"github.com/korableg/V8I.Manager/app/internal/transport/httpserver"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

type (
	Engine struct {
		httpSrvr *httpserver.HttpServer
		sqliteDB *sqlitedb.SqliteDB
		app      *fx.App
	}
)

func NewEngine(cfgPath string) (*Engine, error) {
	logrus.Info("init the app")
	fx.Decorate()
	app := fx.New(
		fx.Provide(validator.New),
		fx.Provide(func(validate *validator.Validate) (config.Config, error) {
			cfg, err := config.NewConfig(cfgPath, validate)
			if err != nil {
				return config.Config{}, fmt.Errorf("init config: %w", err)
			}

			return cfg, nil
		}),
		fx.Provide(sqlitedb.NewSqliteDB),
		user.FX(),
		auth.FX(),
		onecdb.FX(),
		onecserver.FX(),
		client.FX(),
		webinfobase.FX(),
	)

	webCommonInfoBasesHds, err := initWebCommonInfoBasesHandlers(cfg, validate, v8ibuilder)

	httpSrvr := httpserver.NewHttpServer(
		cfg.Http,
		httpserver.WithApiMiddleware(authHds.Middleware()),
		httpserver.WithApiHandlers(userHds),
		httpserver.WithApiHandlers(dbHds),
		httpserver.WithApiHandlers(onecServersHds),
		httpserver.WithHandlers(authHds),
		httpserver.WithHandlers(webCommonInfoBasesHds),
	)

	return &Engine{
		httpSrvr: httpSrvr,
		sqliteDB: sdb,
		app:      app,
	}, nil
}

func (en *Engine) Start() error {
	logrus.Infof("starting http server: %s", en.httpSrvr.Address())

	if err := en.httpSrvr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve http: %w", err)
	}

	return nil
}

func (en *Engine) Shutdown(ctx context.Context) error {
	logrus.Infof("shutdown engine")

	if err := en.httpSrvr.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	if err := en.sqliteDB.Close(); err != nil {
		return fmt.Errorf("close db: %w", err)
	}

	return nil
}

func initDBHandlers(sdb *sqlitedb.SqliteDB, validate *validator.Validate) (*onecdb.Handlers, onecdb.DBCollector, onecdb.V8IBuilder, error) {
	dbRepo, err := onecdb.NewSqliteRepository(sdb)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("init onecdb repository: %w", err)
	}

	dbService, err := onecdb.NewService(dbRepo)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("init onecdb service: %w", err)
	}

	dbHandlers, err := onecdb.NewHandlers(dbService, validate)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("init onecdb handlers: %w", err)
	}

	return dbHandlers, dbService, dbService, nil
}

func initOnecServerHandlers(sdb *sqlitedb.SqliteDB, collector onecdb.DBCollector, validate *validator.Validate) (*onecserver.Handlers, error) {
	onecRepo, err := onecserver.NewSqliteRepository(sdb)
	if err != nil {
		return nil, fmt.Errorf("init onec server repository: %w", err)
	}

	builder, err := dbbuilder.NewBuilder()
	if err != nil {
		return nil, fmt.Errorf("init onec server db builder: %w", err)
	}

	onecService, err := onecserver.NewService(onecRepo, collector, watcher.NewWatcher, builder)
	if err != nil {
		return nil, fmt.Errorf("init onec server service: %w", err)
	}

	onecHandlers, err := onecserver.NewHandlers(onecService, validate)
	if err != nil {
		return nil, fmt.Errorf("init onec server handlers: %w", err)
	}

	return onecHandlers, nil
}

func initWebCommonInfoBasesHandlers(cfg config.Config, validate *validator.Validate, v8iBuilder onecdb.V8IBuilder) (*webinfobase.Handlers, error) {

	svc, err := webinfobase.NewService(cfg.Http.Address, cfg.Http.Port, v8iBuilder)
	if err != nil {
		return nil, fmt.Errorf("init webinfobase service: %w", err)
	}

	webCommonInfoBasesHds, err := webinfobase.NewHandlers(svc, validate)
	if err != nil {
		return nil, fmt.Errorf("init webinfobase handlers: %w", err)
	}

	return webCommonInfoBasesHds, nil

	return nil, nil
}
