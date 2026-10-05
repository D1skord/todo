package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/D1skord/todo/internal/core/domain"
	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
	core_http_types "github.com/D1skord/todo/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"       swaggertype:"string" example:"Спорт"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"Отжиматься 10 раз"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"   swaggertype:"boolean"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")
		}
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchTask 	godoc
// @Summary 	Обновление задачи
// @Description Обновление информации об уже существующей в системе задаче
// @Description ### Логика обновления полей (Three-stage logic):
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"description": "Пойти завтра в спортзал" - устанавливает новое описание задачи в БД`
// @Description 3. **Передан null**: `"description": null` - очищает поле в БД (set to NULL)
// @Description Ограничения: `title` и `completed` не могут быть выставлены как null
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		id 				path int 			true					"ID изменяемой задачи"
// @Param 		request body 	PatchTaskRequest 	true					"PatchTask тело запроса"
// @Success 	200 			{object} PatchTaskResponse 					"Информация о задаче успешно обновлена"
// @Failure 	400 			{object} core_http_response.ErrorResponse 	"Bad Request"
// @Failure 	404 			{object} core_http_response.ErrorResponse 	"Task Not Found"
// @Failure 	409 			{object} core_http_response.ErrorResponse 	"Conflict"
// @Failure 	500 			{object} core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 		/tasks/{id} 	[patch]
func (h *TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
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

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
