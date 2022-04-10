package httpserver

import (
	"go.uber.org/fx"
)

type (
	RouteApiWithAuth struct {
		fx.Out

		RouteRegister RouteRegister `group:"route_api_with_auth"`
	}

	RouteApiAuth struct {
		fx.Out

		RouteRegister RouteRegister `group:"route_api"`
	}

	RouteBase struct {
		fx.Out

		RouteRegister RouteRegister `group:"route_base"`
	}
)
