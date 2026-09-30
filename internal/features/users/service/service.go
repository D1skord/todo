package users_service

import (
	"context"

	"github.com/D1skord/todo/internal/core/domain"
	users_postgres_repository "github.com/D1skord/todo/internal/features/users/repository/postgres"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
	GetUsers(
		ctx context.Context,
		offset *int,
		limit *int,
	) ([]domain.User, error)
	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)
	DeleteUser(
		ctx context.Context,
		id int,
	) error
	PatchUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
}

func NewUserService(usersRepository *users_postgres_repository.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
