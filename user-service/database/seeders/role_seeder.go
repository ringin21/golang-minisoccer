package seeders

import (
	"user-service/domain/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RunRoleSeeder membuat data role awal jika data tersebut belum tersedia di database.
func RunRoleSeeder(db *gorm.DB) {
	// roles berisi daftar role default yang dibutuhkan aplikasi.
	roles := []models.Role{
		{
			Code: "ADMIN",
			Name: "Administrator",
		},
		{
			Code: "Customer",
			Name: "Customer",
		},
	}

	for _, role := range roles {
		// FirstOrCreate mencegah role dibuat berulang jika code role sudah ada.
		err := db.FirstOrCreate(&role, models.Role{Code: role.Code}).Error
		if err != nil {
			logrus.Errorf("failed to seed role: %v", err)
			panic(err)
		}
		logrus.Infof("role %s successfully seeded", role.Code)
	}
}

/*
Kegunaan file:
File ini dibuat untuk menyediakan data role awal aplikasi.
Dengan file ini, role dasar seperti administrator dan customer otomatis tersedia di database.
*/
