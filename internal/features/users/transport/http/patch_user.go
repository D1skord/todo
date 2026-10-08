package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/D1skord/todo/internal/core/domain"
	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
	core_http_types "github.com/D1skord/todo/internal/core/transport/http/types"
)

type PatchUserResponse UserDtoResponse

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"    swaggertype:"string" example:"Ivan Ivanov"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79010050695"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Value == nil {
		return fmt.Errorf(`"full_name" must not be NULL`)
	}

	fullNameLen := len([]rune(*r.FullName.Value))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf(`"full_name" must contain at least 3 characters long`)
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune((*r.PhoneNumber.Value)))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must starts with '+' symbol")
			}
		}
	}

	return nil
}

// PatchUser 	godoc
// @Summary 	Изменение пользователя
// @Description Изменение информации об уже существующем в системе пользователе
// @Description ### Логика обновления полей (Three-stage logic):
// @Description 1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"phone_number": "+79010050695" - устанавливает новый номер телефона в БД`
// @Description 3. **Передан null**: `"phone_number": null` - очищает поле в БД (set to NULL)
// @Description Ограничения: `full_name` не может быть выставлен как null
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		id 				path int 			true					"ID изменяемого пользователя"
// @Param 		request body 	PatchUserRequest 	true					"PatchUser тело запроса"
// @Success 	200 			{object} PatchUserResponse 					"Информация о пользователе успешно изменена"
// @Failure 	400 			{object} core_http_response.ErrorResponse 	"Bad Request"
// @Failure 	404 			{object} core_http_response.ErrorResponse 	"User Not Found"
// @Failure 	500 			{object} core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 		/users/{id} 	[patch]
func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	log.Debug("invoke PatchUser Handler")
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	//Нашли какого пользователя нужно менять
	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get `userID` path value",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	//Нашли как пользователя нужно менять
	userPatch := userPatchFromRequest(request)
	userDomain, err := h.userService.PatchUser(ctx, userId, userPatch)

	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := GetUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(r PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		r.FullName.ToDomain(),
		r.PhoneNumber.ToDomain(),
	)
}
