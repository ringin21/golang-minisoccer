package routes

import (
	"github.com/gin-gonic/gin"
	"user-service/controllers"
	routes "user-service/routes/user"
)

// Registry menyimpan controller registry dan router group utama untuk pendaftaran route.
type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

// IRouteRegister adalah kontrak untuk menjalankan pendaftaran route aplikasi.
type IRouteRegister interface {
	Serve()
}

// NewRouteRegistry membuat instance route registry dengan controller dan router group.
func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRouteRegister {
	return &Registry{controller: controller, group: group}
}

// Serve menjalankan semua route register yang tersedia.
func (r *Registry) Serve() {
	r.userRoute().Run()
}

// userRoute membuat route register khusus untuk fitur user.
func (r *Registry) userRoute() routes.IUserRoute {
	return routes.NewUserRoute(r.controller, r.group)
}

/*
Kegunaan file:
File ini dibuat sebagai pusat pendaftaran route aplikasi.
Dengan file ini, proses registrasi semua route bisa dipanggil dari satu tempat melalui method Serve.
*/
