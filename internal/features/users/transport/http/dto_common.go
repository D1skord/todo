package users_transport_http

import "github.com/D1skord/todo/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func userDTOFromDomain(user domain.User) UserDtoResponse {
	return UserDtoResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(user []domain.User) []UserDtoResponse {
	usersDTO := make([]UserDtoResponse, len(user))

	for i, userDto := range user {
		usersDTO[i] = userDTOFromDomain(userDto)
	}

	return usersDTO
}
