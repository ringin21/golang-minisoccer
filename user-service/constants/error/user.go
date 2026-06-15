package error

import "errors"

var (
	// ErrUserNotFound digunakan ketika data user tidak ditemukan.
	ErrUserNotFound         = errors.New("user not found")
	// ErrPasswordIncorrect digunakan ketika password login tidak sesuai.
	ErrPasswordIncorrect    = errors.New("password incorrect")
	// ErrUsernameExists digunakan ketika username sudah terdaftar.
	ErrUsernameExists       = errors.New("username already exists")
	// ErrPasswordDoesNotMatch digunakan ketika password dan konfirmasi password berbeda.
	ErrPasswordDoesNotMatch = errors.New("password does not match")
)

// UserErrors berisi daftar error terkait proses user yang dikenali oleh aplikasi.
var UserErrors = []error{
	ErrUserNotFound,
	ErrPasswordIncorrect,
	ErrUsernameExists,
	ErrPasswordDoesNotMatch,
}

/*
Kegunaan file:
File ini dibuat untuk menyimpan daftar error yang berhubungan dengan proses user.
Dengan file ini, error seperti user tidak ditemukan, password salah, dan username sudah terdaftar bisa dipakai ulang secara konsisten.
*/
