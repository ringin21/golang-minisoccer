package seeders

import "gorm.io/gorm"

// Registry menyimpan koneksi database yang akan dipakai untuk menjalankan semua seeder.
type Registry struct {
	db *gorm.DB
}

// ISeederRegistry adalah kontrak untuk menjalankan kumpulan seeder.
type ISeederRegistry interface {
	Run()
}

// NewSeederRegistry membuat instance registry seeder dengan koneksi database.
func NewSeederRegistry(db *gorm.DB) ISeederRegistry {
	return &Registry{db: db}
}

// Run menjalankan semua seeder yang dibutuhkan aplikasi.
func (s *Registry) Run() {
	RunRoleSeeder(s.db)
	RunUserSeeder(s.db)
}

/*
Kegunaan file:
File ini dibuat untuk mengatur urutan eksekusi semua seeder.
Dengan file ini, proses seeding cukup dipanggil dari satu registry tanpa memanggil setiap seeder satu per satu.
*/
