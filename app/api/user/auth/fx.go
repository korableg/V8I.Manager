package auth

import "go.uber.org/fx"

func FX() fx.Option {
	return fx.Module("auth",
		fx.Provide(NewAuth),
		fx.Provide(func(a *auth) Auth { return a }),
		fx.Provide(NewHandlers))
}
