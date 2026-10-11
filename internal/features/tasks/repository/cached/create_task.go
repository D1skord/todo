package tasks_cached_repository

import (
	"context"

	"github.com/D1skord/todo/internal/core/domain"
)

func (r *CachedRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	return r.mainRepository.CreateTask(ctx, task)
}
