package services

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"
)

// UserService berisi business logic untuk fitur user.
type UserService struct {
	repository repositories.IRepositoryRegistry
}

// IUserService adalah kontrak method yang disediakan oleh service user.
type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateUserRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

// Claims adalah data yang dimasukkan ke dalam JWT token.
type Claims struct {
	// User berisi informasi user yang akan disimpan di token.
	User *dto.UserResponse
	jwt.RegisteredClaims
}

// NewUserService membuat instance service user dengan repository registry.
func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{repository: repository}
}

// Login memvalidasi username dan password, lalu membuat JWT token jika login berhasil.
func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Cari user berdasarkan username yang dikirim dari request login.
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// Bandingkan password request dengan password hash yang tersimpan di database.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	// Hitung waktu expired token berdasarkan konfigurasi aplikasi.
	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Code),
	}

	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	// Buat dan tanda tangani JWT token memakai secret key dari config.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

// isUsernameExist mengecek apakah username sudah digunakan user lain.
func (u *UserService) isUsernameExist(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByUsername(ctx, username)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

// isEmailExist mengecek apakah email sudah digunakan user lain.
func (u *UserService) isEmailExist(ctx context.Context, email string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

// Register membuat user baru setelah validasi username, email, dan password.
func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Password di-hash sebelum dikirim ke repository untuk disimpan.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Username harus unik agar tidak ada dua user memakai username yang sama.
	if u.isUsernameExist(ctx, req.Username) {
		return nil, errConstant.ErrUsernameExists
	}

	// Email harus unik agar satu email tidak dipakai oleh lebih dari satu user.
	if u.isEmailExist(ctx, req.Email) {
		return nil, errConstant.ErrEmailExists
	}

	// Password dan confirm password harus sama sebelum user dibuat.
	if req.Password != req.ConfirmPassword {
		return nil, errConstant.ErrPasswordDoesNotMatch
	}

	// User baru otomatis diberi role customer.
	user, err := u.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleID:      constants.Customer,
	})
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
		},
	}

	return response, nil
}

// Update memperbarui data user berdasarkan UUID.
func (u *UserService) Update(ctx context.Context, request *dto.UpdateUserRequest, uuid string) (*dto.UserResponse, error) {
	var (
		password                  string
		checkUsername, checkEmail *models.User
		hashedPassword            []byte
		user, userResult          *models.User
		err                       error
		data                      dto.UserResponse
	)

	// Ambil data user lama untuk membandingkan perubahan username dan email.
	user, err = u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Jika username berubah, pastikan username baru belum digunakan user lain.
	isUsernameExist := u.isUsernameExist(ctx, request.Username)
	if isUsernameExist && user.Username != request.Username {
		checkUsername, err = u.repository.GetUser().FindByUsername(ctx, request.Username)
		if err != nil {
			return nil, err
		}

		if checkUsername != nil {
			return nil, errConstant.ErrUsernameExists
		}
	}

	// Jika email berubah, pastikan email baru belum digunakan user lain.
	isEmailExist := u.isEmailExist(ctx, request.Email)
	if isEmailExist && user.Email != request.Email {
		checkEmail, err = u.repository.GetUser().FindByEmail(ctx, request.Email)
		if err != nil {
			return nil, err
		}

		if checkEmail != nil {
			return nil, errConstant.ErrEmailExists
		}
	}

	// Password hanya diubah jika request mengirim password baru.
	if request.Password != "" {
		if request.Password != request.ConfirmPassword {
			return nil, errConstant.ErrPasswordDoesNotMatch
		}
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		password = string(hashedPassword)
	}

	// Kirim data terbaru ke repository untuk di-update ke database.
	userResult, err = u.repository.GetUser().Update(ctx, &dto.UpdateUserRequest{
		Name:        request.Name,
		Username:    request.Username,
		Password:    password,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
	}, uuid)
	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:        userResult.UUID,
		Name:        userResult.Name,
		Username:    userResult.Username,
		PhoneNumber: userResult.PhoneNumber,
		Email:       userResult.Email,
	}

	return &data, nil
}

// GetUserLogin mengambil data user login dari context request.
func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		PhoneNumber: userLogin.PhoneNumber,
		Email:       userLogin.Email,
		Role:        userLogin.Role,
	}

	return &data, nil
}

// GetUserByUUID mengambil data user berdasarkan UUID.
func (u *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	data := dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}

	return &data, nil
}

/*
Kegunaan file:
File ini dibuat untuk menyimpan business logic yang berhubungan dengan user.
Service ini menjadi penghubung antara request/DTO, repository database, validasi password,
validasi data unik, dan pembuatan token autentikasi.
*/
