package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type (
	RouteRegister interface {
		Register(r *mux.Router) *mux.Router
	}

	Config struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	}

	HttpServer struct {
		rootRouter *mux.Router
		apiRouter  *mux.Router
		s          *http.Server
	}

	Option func(server *HttpServer)
)

func WithApiHandlers(r RouteRegister) Option {
	return func(server *HttpServer) {
		r.Register(server.apiRouter)
	}
}

func WithApiMiddleware(m mux.MiddlewareFunc) Option {
	return func(server *HttpServer) {
		server.apiRouter.Use(m)
	}
}

func WithHandlers(r RouteRegister) Option {
	return func(server *HttpServer) {
		r.Register(server.rootRouter)
	}
}

func NewHttpServer(c Config, opts ...Option) *HttpServer {
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	origins := handlers.AllowedOrigins([]string{fmt.Sprintf("http://localhost:%d", c.Port)})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})

	rootRouter := mux.NewRouter()
	rootRouter.Use(handlers.CORS(headers, origins, methods))
	rootRouter.Use(handlers.RecoveryHandler())
	rootRouter.NotFoundHandler = http.NotFoundHandler()

	apiRouter := rootRouter.PathPrefix("/api").Subrouter()

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", c.Address, c.Port),
		Handler: rootRouter,
	}

	h := &HttpServer{
		s:          s,
		rootRouter: rootRouter,
		apiRouter:  apiRouter,
	}

	for _, o := range opts {
		o(h)
	}

	return h
}

func (h *HttpServer) Address() string {
	return h.s.Addr
}

func (h *HttpServer) ListenAndServe() error {
	return h.s.ListenAndServe()
}

func (h *HttpServer) Shutdown(ctx context.Context) error {
	return h.s.Shutdown(ctx)
}
