package helper

import (
	"net/http"

	"github.com/shellrean/extraordinary-raport/domain"
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrInvalidUser:
		return http.StatusNotFound
	case domain.ErrInvalidToken:
		return http.StatusUnauthorized
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}