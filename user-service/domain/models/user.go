package models

import "time"

type User struct {
	ID          uint   `gorm:"primaryKey:autoIncrement"`
	UUID        string `gorm:"type:uuid;not null"`
	Email       string `gorm:"varchar(100);not null"`
	Pass        string `gorm:"varchar(255);not null"`
	Name        string `gorm:"varchar(100);not null"`
	PhoneNumber string `gorm:"varchar(15);not null"`
	RoleID      uint   `gorm:"type: uint;not null"`
	CreateAt    *time.Time
	UpdateAt    *time.Time
	Role        Role `gorm:"foreignKey:role_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
