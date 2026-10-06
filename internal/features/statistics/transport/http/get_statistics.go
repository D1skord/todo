package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_request "github.com/D1skord/todo/internal/core/transport/http/request"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
)

type GetStatisticsResponse StatisticsDTOResponse

// GetStatistics 	godoc
// @Summary 		Получение статистики
// @Description 	Получение статистики по задачам с опциональной фильтрацией по user_id и/или временному промежутку
// @Tags 			statistics
// @Produce 		json
// @Param 			user_id  	query 		int 		false				"Фильтрация стратистики по кнокретному пользователю"
// @Param 			from 		query 		string 		false				"Начало промежутка рассмотрения статистики (включительно), формат: YYYY-MM-DD"
// @Param 			to 			query 		string 		false				"Конец промежутка рассмотрения статистики (не включительно), формат: YYYY-MM-DD"
// @Success 		200 		{object} GetStatisticsResponse 				"Успешное получение статистики"
// @Failure 		400 		{object} core_http_response.ErrorResponse 	"Bad Request"
// @Failure 		500 		{object} core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 			/statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	queryParams, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed get 'userID/from/to' query params",
		)

		return
	}

	statisticsDomain, err := h.statisticsService.GetStatistics(
		ctx,
		queryParams.userID,
		queryParams.from,
		queryParams.to,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get staticstics",
		)

		return
	}

	response := GetStatisticsResponse(toDTOFromDomain(statisticsDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

type queryParams struct {
	userID *int
	from   *time.Time
	to     *time.Time
}

func getUserIDFromToQueryParams(r *http.Request) (*queryParams, error) {
	const (
		toQueryParamKey     = "to"
		fromQueryParamKey   = "from"
		userIDQueryParamKey = "user_id"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetTimeQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetTimeQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return &queryParams{
		userID: userID,
		from:   from,
		to:     to,
	}, nil
}
