package users_transport_http

import (
	"net/http"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
)

type GetUserResponse UserDtoResponse

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
