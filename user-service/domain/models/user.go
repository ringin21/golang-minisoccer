package models

import "time"

// User adalah model database untuk menyimpan data user.
type User struct {
	// ID adalah primary key internal user di database.
	ID          uint   `gorm:"primaryKey:autoIncrement"`
	// UUID adalah identifier unik user yang bisa digunakan di luar database.
	UUID        string `gorm:"type:uuid;not null"`
	// Email berisi alamat email user.
	Email       string `gorm:"varchar(100);not null"`
	// Pass berisi password user yang sudah diproses sebelum disimpan.
	Pass        string `gorm:"varchar(255);not null"`
	// Name berisi nama lengkap user.
	Name        string `gorm:"varchar(100);not null"`
	// PhoneNumber berisi nomor telepon user.
	PhoneNumber string `gorm:"varchar(15);not null"`
	// RoleID adalah foreign key yang menghubungkan user dengan role.
	RoleID      uint   `gorm:"type: uint;not null"`
	// CreateAt berisi waktu data user dibuat.
	CreateAt    *time.Time
	// UpdateAt berisi waktu data user terakhir diperbarui.
	UpdateAt    *time.Time
	// Role berisi relasi user ke tabel role.
	Role        Role `gorm:"foreignKey:role_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

/*
Kegunaan file:
File ini dibuat untuk mendefinisikan model database user.
Model ini digunakan oleh GORM untuk memetakan data user ke tabel database
dan menghubungkan user dengan role melalui relasi foreign key.
*/
