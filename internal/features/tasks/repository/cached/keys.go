package tasks_cached_repository

import "fmt"

func taskKey(id int) string {
	return fmt.Sprintf("task:%d", id)
}
