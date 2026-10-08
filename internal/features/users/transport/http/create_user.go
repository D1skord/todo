package users_transport_http

import (
	"net/http"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79885556677"`
}

type CreateUserResponse UserDtoResponse

// CreateUser 	godoc
// @Summary 	Создать пользователя
// @Description Создать нового пользователя в системе
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		request body 	 CreateUserRequest  true 			"CreateUser тело запроса"
// @Success 	201 	{object} CreateUserResponse 				"Успешно созданный пользователь"
// @Failure 	400 	{object} core_http_response.ErrorResponse 	"Bad Request"
// @Failure 	500 	{object} core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 		/users  [post]
func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	log.Debug("invoke CreateUser Handler")
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.userService.CreateUser(ctx, userDomain)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
