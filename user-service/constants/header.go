package constants

import "net/textproto"

var (
	// XServiceName adalah header untuk mengirim nama service pemanggil.
	XServiceName  = textproto.CanonicalMIMEHeaderKey("x-service-name")
	// XApiKey adalah header untuk mengirim API key antar service.
	XApiKey       = textproto.CanonicalMIMEHeaderKey("x-api-key")
	// XRequestAt adalah header untuk mengirim waktu request dibuat.
	XRequestAt    = textproto.CanonicalMIMEHeaderKey("x-request-at")
	// Authorization adalah header standar untuk token autentikasi.
	Authorization = textproto.CanonicalMIMEHeaderKey("authorization")
)

/*
Kegunaan file:
File ini dibuat untuk menyimpan nama header HTTP yang dipakai aplikasi.
Dengan file ini, penggunaan header menjadi konsisten dan mengurangi risiko salah penulisan nama header.
*/
