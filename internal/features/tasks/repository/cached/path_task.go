package tasks_cached_repository

import (
	"context"

	"github.com/D1skord/todo/internal/core/domain"
)

func (r *CachedRepository) PatchTask(
	ctx context.Context,
	id int,
	task domain.Task,
) (domain.Task, error) {
	patchedTask, err := r.mainRepository.PatchTask(ctx, id, task)
	if err != nil {
		r.deleteTaskFromCache(ctx, id)

		return domain.Task{}, err
	}

	r.cacheTask(ctx, patchedTask)

	return patchedTask, nil
}
