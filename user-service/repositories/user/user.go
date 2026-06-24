package repositories

import (
	"context"
	"errors"
	errWrap "user-service/common/error"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository adalah repository untuk mengakses data user di database.
type UserRepository struct {
	db *gorm.DB
}

// IUserRepository adalah kontrak method yang harus dimiliki repository user.
type IUserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateUserRequest, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
}

// NewUserRepository membuat instance repository user dengan koneksi database.
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

// Register menyimpan data user baru ke database.
func (r *UserRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	// UUID dibuat otomatis agar user memiliki identifier unik yang aman digunakan di luar database.
	user := models.User{
		UUID:        uuid.New().String(),
		Name:        req.Name,
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleID:      req.RoleID,
	}

	// Create menyimpan model user baru menggunakan context request.
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

// Update memperbarui data user berdasarkan UUID.
func (r *UserRepository) Update(ctx context.Context, req *dto.UpdateUserRequest, uuid string) (*models.User, error) {
	user := models.User{
		Name:        req.Name,
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}

	// Where memastikan hanya user dengan UUID tertentu yang diperbarui.
	err := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		Updates(&user).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

// FindByUsername mencari user berdasarkan username.
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	// Preload("Role") digunakan agar data role ikut diambil bersama data user.
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

// FindByEmail mencari user berdasarkan email.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	// Preload("Role") digunakan agar data role ikut diambil bersama data user.
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

// FindByUUID mencari user berdasarkan UUID.
func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User
	// Preload("Role") digunakan agar data role ikut diambil bersama data user.
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("uuid = ?", uuid).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

/*
Kegunaan file:
File ini dibuat untuk menyimpan semua query database yang berhubungan dengan user.
Dengan file ini, akses data user dipusatkan di repository sehingga layer service tidak perlu langsung memakai GORM.
*/
