package dto

// LoginRequest adalah payload request untuk proses login user.
type LoginRequest struct {
	// Username berisi username yang digunakan untuk login.
	Username string `json:"username" validate:"required"`
	// Password berisi password yang digunakan untuk login.
	Password string `json:"password" validate:"required"`
}

// UserResponse adalah format data user yang dikirim sebagai response ke client.
type UserResponse struct {
	// UUID berisi identifier unik user yang aman ditampilkan ke client.
	UUID        string `json:"uuid"`
	// Name berisi nama lengkap user.
	Name        string `json:"name"`
	// Email berisi alamat email user.
	Email       string `json:"email"`
	// Role berisi nama role user.
	Role        string `json:"role"`
	// PhoneNumber berisi nomor telepon user.
	PhoneNumber string `json:"phoneNumber"`
}

// LoginResponse adalah format response setelah login berhasil.
type LoginResponse struct {
	// Token berisi token autentikasi untuk request berikutnya.
	Token string       `json:"token"`
	// User berisi data user yang berhasil login.
	User  UserResponse `json:"user"`
}

// RegisterRequest adalah payload request untuk proses registrasi user.
type RegisterRequest struct {
	// Name berisi nama lengkap user yang akan didaftarkan.
	Name            string `json:"name" validate:"required"`
	// Email berisi email user dan harus memakai format email yang valid.
	Email           string `json:"email" validate:"required,email"`
	// Username berisi username yang akan digunakan untuk login.
	Username        string `json:"username" validate:"required"`
	// Password berisi password user.
	Password        string `json:"password" validate:"required"`
	// ConfirmPassword berisi konfirmasi password untuk memastikan password sesuai.
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	// PhoneNumber berisi nomor telepon user.
	PhoneNumber     string `json:"phoneNumber" validate:"required"`
	// RoleID berisi ID role yang akan diberikan ke user.
	RoleID          uint
}

// RegisterResponse adalah format response setelah registrasi user berhasil.
type RegisterResponse struct {
	// User berisi data user yang berhasil dibuat.
	User UserResponse `json:"user"`
}

// UpdateUserRequest adalah payload request untuk proses update data user.
type UpdateUserRequest struct {
	// Name berisi nama lengkap user terbaru.
	Name            string `json:"name" validate:"required"`
	// Email berisi email user terbaru dan harus memakai format email yang valid.
	Email           string `json:"email" validate:"required,email"`
	// Username berisi username user terbaru.
	Username        string `json:"username" validate:"required"`
	// Password berisi password baru dan bersifat optional.
	Password        string `json:"password,omitempty"`
	// ConfirmPassword berisi konfirmasi password baru dan bersifat optional.
	ConfirmPassword string `json:"confirmPassword,omitempty"`
	// PhoneNumber berisi nomor telepon user terbaru.
	PhoneNumber     string `json:"phoneNumber" validate:"required"`
	// RoleID berisi ID role terbaru untuk user.
	RoleID          uint
}

/*
Kegunaan file:
File ini dibuat untuk menyimpan struktur DTO user.
DTO digunakan sebagai bentuk data yang masuk dari request dan keluar sebagai response,
sehingga layer handler/service tidak langsung memakai model database.
*/
