package tasks_cached_repository

import (
	core_redis_pool "github.com/D1skord/todo/internal/core/repository/redis/pool"
	tasks_service "github.com/D1skord/todo/internal/features/tasks/service"
)

type CachedRepository struct {
	pool           core_redis_pool.Pool
	mainRepository tasks_service.TasksRepository
}

var _ tasks_service.TasksRepository = &CachedRepository{}

func NewCachedRepository(
	pool core_redis_pool.Pool,
	mainRepository tasks_service.TasksRepository,
) *CachedRepository {
	return &CachedRepository{
		pool:           pool,
		mainRepository: mainRepository,
	}
}
