package tasks_cached_repository

import "context"

func (r *CachedRepository) DeleteTask(
	ctx context.Context,
	id int,
) error {
	if err := r.mainRepository.DeleteTask(ctx, id); err != nil {
		return err
	}

	r.deleteTaskFromCache(ctx, id)

	return nil
}
