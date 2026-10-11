package tasks_cached_repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/D1skord/todo/internal/core/domain"
	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_redis_pool "github.com/D1skord/todo/internal/core/repository/redis/pool"
	"go.uber.org/zap"
)

func (r *CachedRepository) getTaskFromCache(
	ctx context.Context,
	id int,
) (domain.Task, bool) {
	log := core_logger.FromContext(ctx)

	bytes, err := r.pool.Get(ctx, taskKey(id))
	if err != nil {
		if errors.Is(err, core_redis_pool.ErrKeyNotFound) {
			log.Debug("task cache miss", zap.Int("task_id", id))
		} else {
			log.Warn("get task from cache", zap.Int("task_id", id), zap.Error(err))
		}

		return domain.Task{}, false
	}

	var model TaskModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		log.Warn("unmarshal task from cache", zap.Int("task_id", id), zap.Error(err))

		return domain.Task{}, false
	}

	log.Debug("task cache hit", zap.Int("task_id", id))

	return taskDomainFromModel(model), true
}

func (r *CachedRepository) cacheTask(
	ctx context.Context,
	task domain.Task,
) {
	log := core_logger.FromContext(ctx)

	bytes, err := json.Marshal(taskModelFromDomain(task))
	if err != nil {
		log.Warn("marshal task for cache", zap.Int("task_id", task.ID), zap.Error(err))

		return
	}

	if err := r.pool.Set(ctx, taskKey(task.ID), bytes, r.pool.TTL()); err != nil {
		log.Warn("set task to cache", zap.Int("task_id", task.ID), zap.Error(err))
	}
}

func (r *CachedRepository) deleteTaskFromCache(
	ctx context.Context,
	id int,
) {
	if err := r.pool.Del(ctx, taskKey(id)); err != nil {
		core_logger.FromContext(ctx).Warn(
			"delete task from cache",
			zap.Int("task_id", id),
			zap.Error(err),
		)
	}
}
