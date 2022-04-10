package onecserver

import (
	"github.com/korableg/V8I.Manager/app/api/onecserver/dbbuilder"
	"github.com/korableg/V8I.Manager/app/api/onecserver/watcher"
	"github.com/korableg/V8I.Manager/app/internal/transport/httpserver"
	"go.uber.org/fx"
)

func FX() fx.Option {
	return fx.Module("onecserver",
		fx.Provide(NewSqliteRepository),
		fx.Provide(func(r *sqliteRepository) Repository { return r }),
		fx.Provide(func() watcher.Fabric { return watcher.NewWatcher }),
		fx.Provide(dbbuilder.NewBuilder()),
		fx.Provide(func(b *dbbuilder.Builder) dbbuilder.DBBuilder { return b }),
		fx.Provide(NewService),
		fx.Provide(func(s *service) Service { return s }),
		fx.Provide(func(h *Handlers) httpserver.RouteApiAuth { return httpserver.RouteApiAuth{RouteRegister: h} }))
}
