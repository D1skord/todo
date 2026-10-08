package users_transport_http

import (
	"net/http"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
)

type GetUserResponse UserDtoResponse

// GetUser 		godoc
// @Summary 	Получение пользователя
// @Description Получение существующего в системе пользователя по его ID
// @Tags 		users
// @Produce 	json
// @Param 		id 			path int true 								"ID получаемого пользователя"
// @Success 	200 		{object} GetUserResponse 					"Пользователь успешно найден"
// @Failure 	400 		{object} core_http_response.ErrorResponse 	"Bad Request"
// @Failure 	404 		{object} core_http_response.ErrorResponse 	"User Not Found"
// @Failure 	500 		{object} core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 		/users/{id} [get]
func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	log.Debug("invoke GetUser Handler")
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get `userID` path value",
		)
		return
	}

	user, err := h.userService.GetUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get `user` path value",
		)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
