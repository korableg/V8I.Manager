package onecdb

import (
	"github.com/korableg/V8I.Manager/app/internal/transport/httpserver"
	"go.uber.org/fx"
)

func FX() fx.Option {
	return fx.Module("onecdb",
		fx.Provide(NewSqliteRepository),
		fx.Provide(func(r *sqliteRepository) Repository { return r }),
		fx.Provide(NewService),
		fx.Provide(func(s *service) (Service, DBCollector, V8IBuilder) { return s, s, s }),
		fx.Provide(NewHandlers),
		fx.Provide(func(h *Handlers) httpserver.RouteApiAuth { return httpserver.RouteApiAuth{RouteRegister: h} }))
}
