package error

import (
	"errors"
	"fmt"
	"strings"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// ValidationResponse adalah format response untuk setiap error validasi field.
type ValidationResponse struct {
	// Field berisi nama field yang gagal divalidasi.
	Field string `json:"field,omitempty"`
	// Message berisi pesan error yang akan dikirim ke client.
	Message string `json:"message,omitempty"`
}

// ErrValidator menyimpan custom message berdasarkan tag validator.
var ErrValidator = map[string]string{}

// ErrValidationResponse mengubah error dari validator menjadi daftar response validasi.
func ErrValidationResponse(err error) (validationResponse []ValidationResponse) {
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				// Tag required digunakan ketika field wajib tidak dikirim.
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is required", err.Field()),
				})
			case "email":
				// Tag email digunakan ketika format email tidak valid.
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a valid email address", err.Field()),
				})
			default:
				// Tag lain akan dicari di ErrValidator untuk custom message.
				errValidator, ok := ErrValidator[err.Tag()]
				if !ok {
					// Jika custom message tidak ditemukan, gunakan pesan fallback.
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("Something wrong on %s; %s", err.Field(), err.Tag()),
					})
					continue
				}

				// Hitung placeholder agar message bisa memakai field saja atau field dan param.
				count := strings.Count(errValidator, "%s")
				if count == 1 {
					// Jika hanya ada satu placeholder, message hanya membutuhkan nama field.
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf(errValidator, err.Field()),
					})
				} else {
					// Jika placeholder lebih dari satu, message membutuhkan nama field dan parameter validator.
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf(errValidator, err.Field(), err.Param()),
					})
				}
			}
		}
	}

	return validationResponse
}

// WrapError mencatat error ke log lalu mengembalikan error yang sama.
func WrapError(err error) error {
	logrus.Errorf("error: %v", err)
	return err
}

/*
Kegunaan file:
File ini dibuat untuk menangani error umum yang dipakai bersama oleh aplikasi.
Isi file ini membantu mengubah error validasi request menjadi response yang mudah dibaca client,
serta menyediakan helper untuk mencatat error ke log.
*/
