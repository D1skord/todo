package core_http_types

import (
	"encoding/json"

	"github.com/D1skord/todo/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

// UnmarshalJSON Вызывается только если пришло значение, либо null
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true // что-то реально передали либо значение, либо null

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
