package config

import (
	"os"
	"user-service/common/util"

	"github.com/sirupsen/logrus"
)

// Config menyimpan konfigurasi aplikasi yang sudah dibaca saat aplikasi berjalan.
var Config AppConfig

// AppConfig adalah struktur utama konfigurasi aplikasi.
type AppConfig struct {
	// Port berisi port yang digunakan aplikasi untuk berjalan.
	Port string `json:"port"`
	// AppName berisi nama aplikasi.
	AppName string `json:"appName"`
	// AppEnv berisi environment aplikasi, misalnya local, development, atau production.
	AppEnv string `json:"appEnv"`
	// SignatureKey berisi key untuk kebutuhan signature antar service.
	SignatureKey string `json:"signatureKey"`
	// Database berisi konfigurasi koneksi database.
	Database Database `json:"database"`
	// RateLimiterMaxRequests berisi jumlah maksimal request dalam periode tertentu.
	RateLimiterMaxRequests float64 `json:"rateLimiterMaxRequests"`
	// RateLimiterTimeSecond berisi durasi rate limiter dalam satuan detik.
	RateLimiterTimeSecond int `json:"rateLimiterTimeSecond"`
	// JwtSecretKey berisi secret key untuk membuat dan memvalidasi JWT.
	JwtSecretKey string `json:"jwtSecretKey"`
	// JwtExpirationTime berisi durasi masa berlaku JWT.
	JwtExpirationTime int `json:"jwtExpirationTime"`
}

// Database adalah struktur konfigurasi koneksi database.
type Database struct {
	// Host berisi host database.
	Host string `json:"host"`
	// Port berisi port database.
	Port int `json:"port"`
	// Name berisi nama database.
	Name string `json:"name"`
	// Username berisi username database.
	Username string `json:"username"`
	// Password berisi password database.
	Password string `json:"password"`
	// MaxOpenConnections berisi jumlah maksimal koneksi database yang dibuka.
	MaxOpenConnections int `json:"maxOpenConnections"`
	// MaxLifeTimeConnection berisi batas waktu hidup koneksi database.
	MaxLifeTimeConnection int `json:"maxLifeTimeConnection"`
	// MaxIdleConnections berisi jumlah maksimal koneksi idle.
	MaxIdleConnections int `json:"maxIdleConnections"`
	// MaxIdleTime berisi batas waktu koneksi berada dalam kondisi idle.
	MaxIdleTime int `json:"maxIdleTime"`
}

// init membaca konfigurasi saat package config pertama kali digunakan.
func init() {
	err := util.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		// Jika config lokal tidak tersedia, aplikasi mencoba membaca config dari Consul.
		logrus.Infof("failed to bind config: %v", err)
		err = util.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}

/*
Kegunaan file:
File ini dibuat untuk menyimpan struktur konfigurasi aplikasi dan proses pembacaan config.
Dengan file ini, aplikasi bisa membaca config dari file JSON lokal terlebih dahulu,
lalu fallback ke Consul jika config lokal tidak tersedia.
*/
