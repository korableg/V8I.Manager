package webinfobase

import "go.uber.org/fx"

func FX() fx.Option {
	return fx.Module("webinfobase",
		fx.Provide(NewService),
		fx.Provide(func(s *service) Service { return s }),
		fx.Provide(NewHandlers),
	)
}
