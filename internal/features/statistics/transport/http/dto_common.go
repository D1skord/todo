package statistics_transport_http

import "github.com/D1skord/todo/internal/core/domain"

type StatisticsDTOResponse struct {
	TasksCreated               int      `json:"tasks_created"                 example:"25"`
	TasksCompleted             int      `json:"tasks_completed"               example:"5"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"          example:"20"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"1m30s"`
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}
	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}
