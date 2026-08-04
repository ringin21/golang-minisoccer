package error

import errField "field-service/constants/error/field"
import errFieldSchedule "field-service/constants/error/fieldschedule"

// ErrMapping mengecek apakah error yang diterima termasuk error yang sudah didefinisikan oleh aplikasi, baik general error maupun user error.
func ErrMapping(err error) bool {
	// Gabungkan daftar error umum dan error terkait user ke dalam satu slice.
	allErrors := make([]error, 0)
	allErrors = append(GeneralErrors[:], errField.FieldErrors[:]...), errFieldSchedule.FieldScheduleErrors[:]...)

	// Bandingkan pesan error input dengan setiap pesan error yang terdaftar.
	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	// Return false jika error tidak ditemukan dalam daftar error aplikasi.
	return false
}

/*
Kegunaan file:
File ini dibuat untuk memetakan error yang sudah didefinisikan oleh aplikasi.
Dengan file ini, aplikasi bisa mengecek apakah error yang terjadi termasuk error yang dikenali atau bukan.
*/
