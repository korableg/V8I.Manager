package user

import "go.uber.org/fx"

func FX() fx.Option {
	return fx.Module("user",
		fx.Provide(NewSqliteRepository),
		fx.Provide(NewService),
		fx.Provide(NewHandlers),
		fx.Provide(func(r *sqliteRepository) Repository { return r }),
		fx.Provide(func(s *service) Service { return s }))
}
