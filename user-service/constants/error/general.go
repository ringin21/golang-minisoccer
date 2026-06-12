package error

import "errors"

const (
	Success = "success"
	Error   = "error"
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSQLError            = errors.New("sql error")
	ErrInvalidToken        = errors.New("invalid token")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrForbidden           = errors.New("forbidden")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrInvalidToken,
	ErrUnauthorized,
	ErrTooManyRequests,
	ErrForbidden,
}
