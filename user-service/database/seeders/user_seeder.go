package seeders

import (
	"user-service/constants"
	"user-service/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RunUserSeeder membuat data user awal jika data tersebut belum tersedia di database.
func RunUserSeeder(db *gorm.DB) {
	// Password default admin di-hash sebelum disimpan ke database.
	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	user := models.User{
		UUID:     uuid.New().String(),
		Name:     "Administrator",
		Username: "admin",
		Password: string(password),
		PhoneNumber: "085151218229",
		Email: "admin@gmail.com",
		RoleID: constants.Admin,
	}

	// FirstOrCreate mencegah data admin dibuat berulang jika username sudah ada.
	err := db.FirstOrCreate(&user, models.User{Username: user.Username}).Error
	if err != nil {
		logrus.Errorf("Failed to seed user: %v", err)
		panic(err)
	}
	logrus.Infof("user %s successfully seeded", user.Username)
}

/*
Kegunaan file:
File ini dibuat untuk menyediakan data user awal aplikasi.
Dengan file ini, aplikasi memiliki akun administrator default yang bisa dipakai setelah database disiapkan.
*/
