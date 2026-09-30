package users_service

import (
	"context"
	"fmt"

	"github.com/D1skord/todo/internal/core/domain"
)

func (s *UsersService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	//1. user.Validate()
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}
	//2. repository level Save(user) - repo is postgres etc.
	user, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
