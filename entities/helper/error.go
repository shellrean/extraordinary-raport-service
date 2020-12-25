package helper

import (
	"github.com/shellrean/extraordinary-raport/domain"
)

func GetErrorCode(err error) int {
	if err == nil {
		return 0
	}
	code, ok := domain.ErrorCode[err]
	if !ok {
		return 1000
	}
	return code
}