package tasks_cached_repository

import (
	"context"

	"github.com/D1skord/todo/internal/core/domain"
)

func (r *CachedRepository) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	if task, ok := r.getTaskFromCache(ctx, id); ok {
		return task, nil
	}

	task, err := r.mainRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}

	r.cacheTask(ctx, task)

	return task, nil
}
