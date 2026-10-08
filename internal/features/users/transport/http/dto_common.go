package users_transport_http

import "github.com/D1skord/todo/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id" example:"1"`
	Version     int     `json:"version" example:"1"`
	FullName    string  `json:"full_name" example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" example:"+79885556677"`
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
