package controllers

import (
	"net/http"
	errWrap "user-service/common/error"
	"user-service/common/response"
	"user-service/domain/dto"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserController menangani request HTTP yang berhubungan dengan user.
type UserController struct {
	service services.IServiceRegistry
}

// IUserController adalah kontrak endpoint yang disediakan controller user.
type IUserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Update(*gin.Context)
	GetUserLogin(*gin.Context)
	GetUserByUUID(*gin.Context)
}

// NewUserController membuat instance controller user dengan service registry.
func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{service: service}
}

// Login menerima request login, memvalidasi payload, lalu mengembalikan token jika berhasil.
func (u *UserController) Login(ctx *gin.Context) {
	request := &dto.LoginRequest{}

	// ShouldBindJSON membaca body JSON request ke struct LoginRequest.
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Validator mengecek aturan validate tag pada DTO.
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	// Jika payload valid, proses login diteruskan ke service user.
	user, err := u.service.GetUser().Login(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code:  http.StatusOK,
		Data:  user.User,
		Token: &user.Token,
		Gin:   ctx,
	})
}

// Register menerima request registrasi, memvalidasi payload, lalu membuat user baru.
func (u *UserController) Register(ctx *gin.Context) {
	request := &dto.RegisterRequest{}

	// ShouldBindJSON membaca body JSON request ke struct RegisterRequest.
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Validator mengecek aturan validate tag pada DTO.
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	// Jika payload valid, proses registrasi diteruskan ke service user.
	user, err := u.service.GetUser().Register(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user.User,
		Gin:  ctx,
	})
}

// Update menerima request update user berdasarkan UUID dari path parameter.
func (u *UserController) Update(ctx *gin.Context) {
	request := &dto.UpdateUserRequest{}
	// UUID diambil dari parameter route.
	uuid := ctx.Param("uuid")

	// ShouldBindJSON membaca body JSON request ke struct UpdateUserRequest.
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Validator mengecek aturan validate tag pada DTO.
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	// Jika payload valid, proses update diteruskan ke service user.
	user, err := u.service.GetUser().Update(ctx, request, uuid)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

// GetUserLogin mengambil data user yang sedang login dari context request.
func (u *UserController) GetUserLogin(ctx *gin.Context) {
	// Service membaca data user login dari context yang sudah diisi middleware.
	user, err := u.service.GetUser().GetUserLogin(ctx.Request.Context())
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

// GetUserByUUID mengambil data user berdasarkan UUID dari path parameter.
func (u *UserController) GetUserByUUID(ctx *gin.Context) {
	// UUID diambil dari parameter route dan diteruskan ke service user.
	user, err := u.service.GetUser().GetUserByUUID(ctx.Request.Context(), ctx.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

/*
Kegunaan file:
File ini dibuat untuk menangani request HTTP terkait user.
Controller bertugas membaca input dari Gin context, melakukan validasi request,
memanggil service user, lalu mengirim response JSON standar ke client.
*/
