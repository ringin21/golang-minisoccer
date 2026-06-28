package cmd

import (
	"fmt"
	"net/http"
	"time"
	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	"user-service/controllers"
	"user-service/database/seeders"
	"user-service/domain/models"
	"user-service/middlewares"
	"user-service/repositories"
	"user-service/routes"
	"user-service/services"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// command adalah perintah CLI untuk menjalankan HTTP server user-service.
var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(c *cobra.Command, args []string) {
		// Load file .env agar environment variable lokal bisa dibaca saat aplikasi dijalankan.
		_ = godotenv.Load()
		// Inisialisasi koneksi database menggunakan konfigurasi aplikasi.
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		// Set timezone aplikasi ke Asia/Jakarta agar waktu default mengikuti timezone service.
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		// AutoMigrate membuat atau menyesuaikan tabel database berdasarkan model GORM.
		err = db.AutoMigrate(
			&models.Role{},
			&models.User{},
		)
		if err != nil {
			panic(err)
		}

		// Seeder dijalankan untuk memastikan data awal seperti role dan admin tersedia.
		seeders.NewSeederRegistry(db).Run()
		// Registry dibuat berurutan dari repository, service, sampai controller sebagai dependency aplikasi.
		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository)
		controller := controllers.NewControllerRegistry(service)

		// Gin router utama dibuat untuk menerima dan memproses request HTTP.
		router := gin.Default()
		// Middleware panic recovery dipasang agar panic tidak langsung menjatuhkan server.
		router.Use(middlewares.HandlePanic())
		// NoRoute menangani request ke path yang tidak terdaftar.
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})
		// Root endpoint digunakan sebagai health/welcome endpoint sederhana.
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to User Service",
			})
		})
		// Middleware CORS mengatur origin, method, dan header yang diizinkan oleh API.
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-request-at, x-api-key")
			// Request OPTIONS adalah preflight CORS, jadi cukup dikembalikan status 204.
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		// Rate limiter dibuat dari konfigurasi untuk membatasi jumlah request dalam durasi tertentu.
		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequests,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			})
		// Middleware rate limiter dipasang ke seluruh route.
		router.Use(middlewares.RateLimiter(lmt))

		// Seluruh route API versi 1 memakai prefix /api/v1.
		group := router.Group("/api/v1")
		// Route registry mendaftarkan semua route aplikasi ke router group.
		route := routes.NewRouteRegistry(controller, group)
		route.Serve()

		// Port server dibaca dari konfigurasi dan dijalankan dengan format ":port".
		port := fmt.Sprintf(":%s", config.Config.Port)
		router.Run(port)
	},
}

// Run menjalankan command Cobra untuk memulai aplikasi.
func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}

/*
Kegunaan file:
File ini dibuat sebagai entrypoint command untuk menjalankan user-service.
Di file ini aplikasi melakukan bootstrap utama seperti load environment, koneksi database,
migration, seeding data awal, pembuatan dependency registry, setup middleware,
registrasi route, dan menjalankan HTTP server.
*/
