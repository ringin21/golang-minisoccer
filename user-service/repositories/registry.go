package repositories

import (
	repositories "user-service/repositories/user"

	"gorm.io/gorm"
)

// Registry menyimpan koneksi database yang dipakai untuk membuat repository.
type Registry struct {
	db *gorm.DB
}

// IRepositoryRegistry adalah kontrak untuk mengambil repository yang tersedia.
type IRepositoryRegistry interface {
	GetUser() repositories.IUserRepository
}

// NewRepositoryRegistry membuat instance registry repository dengan koneksi database.
func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

// GetUser mengembalikan repository user yang siap digunakan.
func (r *Registry) GetUser() repositories.IUserRepository {
	return repositories.NewUserRepository(r.db)
}

/*
Kegunaan file:
File ini dibuat sebagai pusat pendaftaran repository.
Dengan file ini, layer lain cukup mengambil repository melalui registry tanpa membuat repository secara langsung.
*/
