package users_transport_http

import (
	"net/http"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
	core_http_utils "github.com/D1skord/todo/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	log.Debug("invoke DeleteUser Handler")
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
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
