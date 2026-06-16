package response

import (
	"net/http"
	"user-service/constants"
	errConstant "user-service/constants/error"

	"github.com/gin-gonic/gin"
)

// Response adalah format standar body JSON yang dikirim ke client.
type Response struct {
	// Status berisi status response, bisa "success" atau "error".
	Status string `json:"status"`
	// Message berisi pesan yang akan dikirim ke client, biasanya digunakan untuk error.
	Message string `json:"message"`
	// Data berisi data yang akan dikirim ke client, biasanya digunakan untuk response sukses.
	Data interface{} `json:"data"`
	// Token berisi token autentikasi yang akan dikirim ke client.
	Token string `json:"token,omitempty"`
}

// ParameterHTTPResp berisi semua parameter yang dibutuhkan untuk membuat HTTP response.
type ParameterHTTPResp struct {
	// Code adalah HTTP status code yang akan dikirim ke client.
	Code int
	// Error berisi error yang menentukan response sukses atau gagal.
	Error error
	// Message memakai *string agar bisa bernilai nil ketika caller tidak mengirim custom message.
	Message *string
	// Gin memakai *gin.Context karena Gin mengelola request/response melalui context pointer.
	Gin *gin.Context
	// Data berisi payload response dan bisa menerima tipe data apa pun.
	Data interface{}
	// Token memakai *string agar token bersifat optional dan bisa nil jika tidak perlu dikirim.
	Token *string
}

// HTTPResponse membuat response JSON standar berdasarkan parameter yang diterima.
func HTTPResponse(param ParameterHTTPResp) {
	token := ""
	if param.Token != nil {
		token = *param.Token
	}

	// Jika tidak ada error, kirim response sukses.
	if param.Error == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
			Token:   token,
		})
		return
	}

	// Jika ada error, gunakan internal server error sebagai pesan default.
	message := errConstant.ErrInternalServerError.Error()
	if param.Message != nil {
		message = *param.Message
	} else if param.Error != nil {
		// Jika error termasuk error aplikasi yang dikenali, gunakan pesan error tersebut.
		if errConstant.ErrMapping(param.Error) {
			message = param.Error.Error()
		}
	}

	// Kirim response error dengan format JSON standar.
	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    param.Data,
	})
}

/*
Kegunaan file:
File ini dibuat untuk menyatukan format response HTTP yang dikirim oleh aplikasi.
Dengan file ini, setiap endpoint bisa mengirim response sukses atau error dengan struktur JSON yang konsisten.
*/
