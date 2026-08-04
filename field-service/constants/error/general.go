package error

import "errors"

var (
	// ErrInternalServerError digunakan untuk error umum yang tidak terduga dari server.
	ErrInternalServerError = errors.New("internal server error")
	// ErrSQLError digunakan ketika terjadi error pada proses database atau query SQL.
	ErrSQLError            = errors.New("sql error")
	// ErrInvalidToken digunakan ketika token yang diterima tidak valid.
	ErrInvalidToken        = errors.New("invalid token")
	// ErrUnauthorized digunakan ketika request tidak memiliki akses autentikasi yang valid.
	ErrUnauthorized        = errors.New("unauthorized")
	// ErrTooManyRequests digunakan ketika request melebihi batas rate limit.
	ErrTooManyRequests     = errors.New("too many requests")
	// ErrForbidden digunakan ketika user terautentikasi tetapi tidak memiliki izin akses.
	ErrForbidden           = errors.New("forbidden")
)

// GeneralErrors berisi daftar error umum yang dikenali oleh aplikasi.
var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrInvalidToken,
	ErrUnauthorized,
	ErrTooManyRequests,
	ErrForbidden,
}

/*
Kegunaan file:
File ini dibuat untuk menyimpan status response dan daftar general error aplikasi.
Dengan file ini, error umum seperti unauthorized, forbidden, dan internal server error bisa dipakai ulang secara konsisten.
*/
