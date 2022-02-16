package httpserver

import (
	"github.com/gorilla/mux"
)

type (
	RouteRegister interface {
		Register(r *mux.Router) *mux.Router
	}

	HttpServer struct {
	}
)

func (h *HttpServer) Start() error {

}
