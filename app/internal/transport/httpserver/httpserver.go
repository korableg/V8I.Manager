package httpserver

import (
	"github.com/gorilla/mux"
)

type Router interface {
	Register(r *mux.Router) error
}
