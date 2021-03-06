package api

import (
	"net/http"

	"github.com/shellrean/extraordinary-raport/domain"
)

type response map[string]interface{}

type ErrorValidation struct {
	Field		string 		`json:"field"`
	Message 	string 		`json:"message"`
}

var errorCodeToResponseCode = map[error]int {
	domain.ErrServerError: 		http.StatusInternalServerError,
	domain.ErrHeaderMiss: 		http.StatusBadRequest,
	domain.ErrParamMiss: 		http.StatusBadRequest,
	domain.ErrInvalidOffset:	http.StatusBadRequest,
	domain.ErrInvalidLocale: 	http.StatusBadRequest,
	domain.ErrInvalidTimezone: 	http.StatusBadRequest,
	domain.ErrTooManyRequest: 	http.StatusTooManyRequests,
	domain.ErrBadParamInput:	http.StatusBadRequest,
	domain.ErrExistData:		http.StatusBadRequest,
	domain.ErrNotExistData:		http.StatusBadRequest,
	domain.ErrUnauthorized: 	http.StatusUnauthorized,
	domain.ErrNoAuthorized: 	http.StatusForbidden,
	domain.ErrUnprocess: 		http.StatusUnprocessableEntity,
	domain.ErrAuthFailed: 		http.StatusUnauthorized,
	domain.ErrNotFound: 		http.StatusNotFound,
	domain.ErrValidation: 		http.StatusLengthRequired,
	domain.ErrFileNotAllowed:	http.StatusBadRequest,
	domain.ErrSessExpired: 		http.StatusUnauthorized,
	domain.ErrSessVerifation: 	http.StatusUnauthorized,
	domain.ErrSessDecode: 		http.StatusUnauthorized,
	domain.ErrSessInvalid: 		http.StatusUnauthorized,
	domain.ErrUnauthorizedUser: http.StatusUnauthorized,
	domain.ErrUserNotFound: 	http.StatusUnauthorized,
	domain.ErrCredential: 		http.StatusUnauthorized,
	domain.ErrLoginTypeInvalid: http.StatusUnauthorized,
	domain.ErrSocialTypeInvalid:http.StatusUnauthorized,
	domain.ErrLogin: 			http.StatusUnauthorized,
	domain.ErrAccountDisable: 	http.StatusUnauthorized,
	domain.ErrClassroomAcademicNotFound: http.StatusNotFound,
	domain.ErrSubjectNotFound:	http.StatusNotFound,
	domain.ErrUserDataNotFound:	http.StatusNotFound,
	domain.ErrClassroomSubjectNotFound: http.StatusNotFound,
	domain.ErrClassroomNotFound: http.StatusNotFound,
	domain.ErrSettingNotFound: 	http.StatusNotFound,
	domain.ErrClassroomStudentNotFound: http.StatusNotFound,
	domain.ErrStudentNotFound:	http.StatusNotFound,
	domain.ErrAcademicNotFound:	http.StatusNotFound,
	domain.ErrSubjectPlanNotFound: http.StatusNotFound,
	domain.ErrExschoolNotFound:	http.StatusNotFound,
	domain.ErrExschoolStudentNotFound: http.StatusNotFound,
	domain.ErrNoteTypeNotFound: http.StatusNotFound,
	domain.ErrAcademicYearExist: http.StatusBadRequest,
	domain.ErrClassroomAcademicExist: http.StatusBadRequest,
	domain.ErrClassroomStudentExist: http.StatusBadRequest,
	domain.ErrEmailExist:		http.StatusBadRequest,
}

func ResponseSuccess(msg string, data interface{}) response{
	return response{
		"success":	true,
		"message": 	msg,
		"data":		data,
	}
}

func ResponseError(msg string, code int) response{
	return response{
		"success":		false,
		"message":		msg,
		"error_code":	code,
	}
}

func ResponseErrorWithData(msg string, code int, data interface{}) response{
	return response{
		"success":		false,
		"message":		msg,
		"errors":		data,
		"error_code":	code,
	}
}

func GetHttpStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	code, ok := errorCodeToResponseCode[err]
	if !ok {
		return http.StatusInternalServerError
	}
	return code
}