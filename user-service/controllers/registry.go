package controllers

import (
	controllers "user-service/controllers/user"
	"user-service/services"
)

// Registry menyimpan dependency service yang dibutuhkan untuk membuat controller.
type Registry struct {
	service services.IServiceRegistry
}

// IControllerRegistry adalah kontrak untuk mengambil controller yang tersedia.
type IControllerRegistry interface {
	GetUserController() controllers.IUserController
}

// NewControllerRegistry membuat instance registry controller dengan service registry.
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

// GetUserController mengembalikan controller user yang siap digunakan.
func (u *Registry) GetUserController() controllers.IUserController {
	return controllers.NewUserController(u.service)
}

/*
Kegunaan file:
File ini dibuat sebagai pusat pendaftaran controller.
Dengan file ini, router atau layer lain cukup mengambil controller melalui registry
tanpa membuat controller secara langsung.
*/
