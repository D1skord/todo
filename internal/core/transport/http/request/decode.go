package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/D1skord/todo/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	var err error

	v, ok := dest.(validatable)
	//Если у кастомной структуры есть кастомный метод валидации
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"validate request: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
