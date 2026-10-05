package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
)

// DeleteTask 	godoc
// @Summary 	Удалить задачу
// @Description Удалить существующего в системе задачу по его ID
// @Tags 		tasks
// @Param 		id path 	int true 									"ID удаляемой задачи"
// @Success 	204 													"Успешное удаление задачи"
// @Failure 	400 		{object} core_http_response.ErrorResponse 	"Bad Request"
// @Failure 	404 		{object} core_http_response.ErrorResponse 	"Task  Not Found"
// @Failure 	500 		{object} core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 		/tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)

		return
	}

	responseHandler.NoContentResponse()
}
