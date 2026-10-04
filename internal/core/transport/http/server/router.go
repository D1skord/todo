package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/D1skord/todo/internal/core/transport/http/middleware"
)

type APIVersion string

var (
	APIVersion1 = APIVersion("v1")
	APIVersion2 = APIVersion("v2")
	APIVersion3 = APIVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
	middleware []core_http_middleware.Middleware
}

func NewAPIVersionRouter(
	version APIVersion,
	middleware ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: version,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		// "GET /tasks"
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
