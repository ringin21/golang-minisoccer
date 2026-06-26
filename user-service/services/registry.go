package services

import (
	"user-service/repositories"
	services "user-service/services/user"
)

// Registry menyimpan dependency repository yang dibutuhkan untuk membuat service.
type Registry struct {
	repository repositories.IRepositoryRegistry
}

// IServiceRegistry adalah kontrak untuk mengambil service yang tersedia.
type IServiceRegistry interface {
	GetUser() services.IUserService
}

// NewServiceRegistry membuat instance registry service dengan repository registry.
func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

// GetUser mengembalikan service user yang siap digunakan.
func (r *Registry) GetUser() services.IUserService {
	return services.NewUserService(r.repository)
}

/*
Kegunaan file:
File ini dibuat sebagai pusat pendaftaran service.
Dengan file ini, layer lain cukup mengambil service melalui registry tanpa membuat service secara langsung.
*/
