package constant

import (
	"net/http"

	"github.com/AriaPutra01/go-commerce/internal/exception"
)

var (
	ErrInvalidCredentials  = &exception.AppError{StatusCode: http.StatusUnauthorized, ErrCode: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrEmailAlreadyExists  = &exception.AppError{StatusCode: http.StatusConflict, ErrCode: "EMAIL_ALREADY_EXISTS", Message: "Email has been registered"}
	ErrInvalidRefreshToken = &exception.AppError{StatusCode: http.StatusUnauthorized, ErrCode: "INVALID_REFRESH_TOKEN", Message: "The login session has expired, please log in again."}
)
