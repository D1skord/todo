package tasks_cached_repository

import (
	"context"

	"github.com/D1skord/todo/internal/core/domain"
)

func (r *CachedRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	return r.mainRepository.GetTasks(ctx, userID, limit, offset)
}
