package routes

import (
	"github.com/gin-gonic/gin"
	"user-service/controllers"
	"user-service/middlewares"
)

// UserRoute menyimpan controller dan router group untuk route user.
type UserRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

// IUserRoute adalah kontrak untuk menjalankan pendaftaran route user.
type IUserRoute interface {
	Run()
}

// NewUserRoute membuat instance route user dengan controller registry dan router group.
func NewUserRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IUserRoute {
	return &UserRoute{controller: controller, group: group}
}

// Run mendaftarkan semua endpoint yang berhubungan dengan user dan autentikasi.
func (u *UserRoute) Run() {
	// Group /auth menjadi prefix untuk seluruh route user/authentication.
	group := u.group.Group("/auth")
	// GET /auth/user mengambil data user yang sedang login dan membutuhkan autentikasi.
	group.GET("/user", middlewares.Authenticate(), u.controller.GetUserController().GetUserLogin)
	// GET /auth/:uuid mengambil data user berdasarkan UUID dan membutuhkan autentikasi.
	group.GET("/:uuid", middlewares.Authenticate(), u.controller.GetUserController().GetUserByUUID)
	// POST /auth/login digunakan untuk login dan mendapatkan token.
	group.POST("/login", u.controller.GetUserController().Login)
	// POST /auth/register digunakan untuk registrasi user baru.
	group.POST("/register", u.controller.GetUserController().Register)
	// PUT /auth/:uuid digunakan untuk update data user dan membutuhkan autentikasi.
	group.PUT("/:uuid", middlewares.Authenticate(), u.controller.GetUserController().Update)
}

/*
Kegunaan file:
File ini dibuat untuk mendaftarkan endpoint yang berhubungan dengan user dan autentikasi.
Dengan file ini, path route, controller handler, dan middleware authentication dikelompokkan dalam satu tempat.
*/
