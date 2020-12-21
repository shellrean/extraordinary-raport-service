package helper

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
)

func GetErrorMessage(e validator.FieldError) (message string) {
	switch e.Tag() {
	case "required":
		message = fmt.Sprintf("%s must be filled", e.Field())	
	case "email":
		message = fmt.Sprintf("%s must be valid email", e.Field())
	case "gte":
		message = fmt.Sprintf("%s must be greater or equal than %s", e.Field(), e.Param())
	case "lte":
		message = fmt.Sprintf("%s must be less or equal than %s", e.Field(), e.Param())
	}
	
	return
}