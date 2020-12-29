package domain

import "errors"

// error code in the system
// to ensure that our code recognized by developer
var (
	// (1000) internal app server process error
	ErrServerError = errors.New("App Server Error, please contact the admin")
	
	// (1001) missing header request
	ErrHeaderMiss = errors.New("Missing Headers")

	// (1002) missing parameter
	ErrParamMiss = errors.New("Missing Parameters")

	// (1003) invalid offset of limit
	ErrInvalidOffset = errors.New("Invalid offset or limit")

	// (1004) invalid locale
	ErrInvalidLocale = errors.New("Invalid Locale")

	// (1005) invalid timezone
	ErrInvalidTimezone = errors.New("Invalid Timezone")

	// (1006) exceeded request per minute
	ErrTooManyRequest = errors.New("You exceeded the limit of requests per minute, Please try again after sometime")
	
	// (1007) bad parameter input
	ErrBadParamInput = errors.New("Bad parameter input")

	// (1008) data is existing
	ErrExistData = errors.New("Data is exist")

	// (1101) unauthorized
	ErrUnauthorized = errors.New("Unauthorized")

	// (1102) not authorized to access
	ErrNoAuthorized = errors.New("Not authorized to access")

	// (1103) unprocesable entity
	ErrUnprocess = errors.New("Unprocessable Entity")

	// (1104) authentication failed
	ErrAuthFailed = errors.New("Authentication Failed")

	// (1105) not found
	ErrNotFound = errors.New("Not Found")

	// (1106) validation error
	ErrValidation = errors.New("Validation error")

	// (1201) session expired
	ErrSessExpired = errors.New("Your session is expired, please login again")

	// (1202) session verification error
	ErrSessVerifation = errors.New("Your sessions is invalid")

	// (1203) session decoding error
	ErrSessDecode = errors.New("Your session sessions is invalid")

	// (1204) invalid session
	ErrSessInvalid = errors.New("Your sessions token is invalid")

	// (1205) unauthorized user
	ErrUnauthorizedUser = errors.New("You are Unauthorized, Please login")

	// (1206) user not found
	ErrUserNotFound = errors.New("Authentication Error, User Not found")

	// (1301) invalid credentials
	ErrCredential = errors.New("Invalid Credentials")

	// (1302) invalid login type
	ErrLoginTypeInvalid = errors.New("Invalid Login Type")

	// (1303) invalid social type
	ErrSocialTypeInvalid = errors.New("Invalid Social Type")

	// (1304) login error
	ErrLogin = errors.New("Login Error")

	// (1305) account disabled
	ErrAccountDisable = errors.New("Your Account is disabled by the admin")
)

var ErrorCode = map[error]int{
	ErrServerError: 	1000,
	ErrHeaderMiss:		1001,
	ErrParamMiss:		1002,
	ErrInvalidOffset:	1003,
	ErrInvalidLocale:	1004,
	ErrInvalidTimezone:	1005,
	ErrTooManyRequest:	1006,
	ErrBadParamInput:	1007,
	ErrExistData:		1008,
	ErrUnauthorized:	1101,
	ErrNoAuthorized:	1102,
	ErrUnprocess:		1103,
	ErrAuthFailed:		1104,
	ErrNotFound:		1105,
	ErrValidation:		1106,
	ErrSessExpired:		1201,
	ErrSessVerifation:	1202,
	ErrSessDecode:		1203,
	ErrSessInvalid:		1204,
	ErrUnauthorizedUser:1205,
	ErrUserNotFound:	1206,
	ErrCredential:		1301,
	ErrLoginTypeInvalid: 1302,
	ErrSocialTypeInvalid: 1303,
	ErrLogin:			1304,
	ErrAccountDisable:	1305,
}