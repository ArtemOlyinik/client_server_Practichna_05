package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func validationErrors(err error) map[string]string {
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		result := make(map[string]string, len(fieldErrors))
		for _, fieldError := range fieldErrors {
			result[fieldError.Field()] = fieldError.Tag()
		}
		return result
	}
	return map[string]string{"body": err.Error()}
}
