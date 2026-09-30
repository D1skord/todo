package users_transport_http

import (
	"context"
	"net/http"

	"github.com/D1skord/todo/internal/core/domain"
	core_http_server "github.com/D1skord/todo/internal/core/transport/http/server"
	users_service "github.com/D1skord/todo/internal/features/users/service"
)

type UsersHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, patch domain.UserPatch) (domain.User, error)
}

func NewUsersHTTPHandler(userService *users_service.UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		userService: userService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			http.MethodPost,
			"/users",
			h.CreateUser,
		},
		{
			http.MethodGet,
			"/users",
			h.GetUsers,
		},
		{
			http.MethodGet,
			"/users/{id}",
			h.GetUser,
		}, {
			http.MethodDelete,
			"/users/{id}",
			h.DeleteUser,
		}, {
			http.MethodPatch,
			"/users/{id}",
			h.PatchUser,
		},
	}
}
