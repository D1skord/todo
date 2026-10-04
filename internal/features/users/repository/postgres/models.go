package users_postgres_repository

import "github.com/D1skord/todo/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userModels := make([]domain.User, len(users))

	for i, user := range users {
		userModels[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
		)
	}

	return userModels
}
