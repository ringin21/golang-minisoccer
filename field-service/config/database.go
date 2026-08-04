package config

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDatabase membuat koneksi database PostgreSQL menggunakan konfigurasi aplikasi.
func InitDatabase() (*gorm.DB, error) {
	config := Config
	// Password di-escape agar karakter khusus tetap aman saat dimasukkan ke URI database.
	endcodedPassword := url.QueryEscape(config.Database.Password)
	// URI digunakan oleh driver PostgreSQL untuk membuka koneksi database.
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.Database.Username, endcodedPassword, config.Database.Host, config.Database.Port, config.Database.Name)

	// Membuka koneksi database menggunakan GORM dan driver PostgreSQL.
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Mengambil instance database SQL bawaan untuk mengatur connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Mengatur batas koneksi database berdasarkan konfigurasi aplikasi.
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeTimeConnection) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdleTime) * time.Second)
	return db, nil
}

/*
Kegunaan file:
File ini dibuat untuk mengatur proses inisialisasi koneksi database.
Dengan file ini, aplikasi bisa membuat koneksi PostgreSQL melalui GORM dan mengatur connection pool
berdasarkan konfigurasi yang sudah dibaca dari config aplikasi.
*/
