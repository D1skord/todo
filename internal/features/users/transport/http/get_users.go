package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/D1skord/todo/internal/core/domain"
	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
)

type GetUsersResponse []UserDtoResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	log.Debug("invoke GetUsers Handler")
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)

	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failder to get `limit` query parameter",
		)
		return
	}

	userDomains, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)

		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)

	if err != nil {
		return nil, nil, fmt.Errorf("get `%s` parameter: %w", limitQueryParamKey, err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get `%s` parameter: %w", offsetQueryParamKey, err)
	}

	return limit, offset, nil
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
