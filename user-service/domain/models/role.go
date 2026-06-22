package models

import "time"

// Role adalah model database untuk menyimpan data role user.
type Role struct {
	// ID adalah primary key internal role di database.
	ID        uint    `gorm:"primaryKey:autoIncrement"`
	// Code berisi kode role yang digunakan aplikasi.
	Code      string  `gorm:"varchar(15);not null"`
	// Name berisi nama role yang mudah dibaca.
	Name      string  `gorm:"varchar(20);not null"`
	// CreatedAt berisi waktu data role dibuat.
	CreatedAt *time.Time 
	// UpdatedAt berisi waktu data role terakhir diperbarui.
	UpdatedAt *time.Time
}

/*
Kegunaan file:
File ini dibuat untuk mendefinisikan model database role.
Model ini digunakan oleh GORM untuk memetakan data role ke tabel database
dan menjadi referensi relasi untuk model user.
*/
