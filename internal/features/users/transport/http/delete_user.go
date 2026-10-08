package users_transport_http

import (
	"net/http"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
)

// DeleteUser 	godoc
// @Summary 	Удалить пользователя
// @Description Удалить существующего в системе пользователя по его ID
// @Tags 		users
// @Param 		id path 	int true 									"ID удаляемого пользователя"
// @Success 	204 													"Успешное удаление пользователя"
// @Failure 	400 		{object} core_http_response.ErrorResponse 	"Bad Request"
// @Failure 	404 		{object} core_http_response.ErrorResponse 	"User Not Found"
// @Failure 	500 		{object} core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 		/users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	log.Debug("invoke DeleteUser Handler")
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get `userID` path value",
		)
		return
	}

	if err = h.userService.DeleteUser(ctx, userId); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete `userId` path value",
		)
		return
	}

	responseHandler.NoContentResponse()
}
