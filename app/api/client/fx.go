package client

import "go.uber.org/fx"

func FX() fx.Option {
	return fx.Module("client",
		fx.Provide(NewService),
		fx.Provide(func(s *service) Service { return s }))
}
