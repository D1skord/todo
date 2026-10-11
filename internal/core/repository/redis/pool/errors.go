package core_redis_pool

import "errors"

var (
	ErrKeyNotFound = errors.New("redis: key not found")
)
