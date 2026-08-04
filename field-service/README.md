# Go Skeleton - Microservice Template

Template boilerplate untuk membuat microservice menggunakan Go dengan pola arsitektur Clean Architecture.

## Struktur Folder dan File

### File Root

| File | Kegunaan |
|------|----------|
| `main.go` | Titik masuk (entry point) utama aplikasi Go. Tempat menjalankan fungsi `main()` yang menginisialisasi dan menjalankan server. |
| `go.mod` | File modul Go yang mendefinisikan nama module (`user-service`) dan dependency management proyek. |
| `Makefile` | Kumpulan perintah (command) untuk memudahkan operasi development seperti build, run, hot reload, dan Docker operations. |
| `Dockerfile` | Instruksi untuk membangun Docker image agar aplikasi bisa dijalankan dalam container. |
| `docker-compose.yml` | Konfigurasi untuk menjalankan layanan (service) beserta dependensinya menggunakan Docker Compose. |
| `Jenkinsfile` | Pipeline otomasi CI/CD untuk integrasi dengan Jenkins. |
| `.gitignore` | Daftar file dan folder yang diabaikan oleh Git (tidak di-track ke dalam repository). |

---

### Folder `cmd/`

Menyimpan kode untuk **entry point** atau **command execution**. Biasanya berisi file-file yang menjalankan aplikasi seperti `main.go` atau sub-command untuk CLI tools. Folder ini memisahkan titik masuk aplikasi dari logika bisnis.

---

### Folder `config/`

Berisi kode untuk **mengelola konfigurasi aplikasi** seperti environment variables, database connection, Redis, dan konfigurasi lainnya. Konfigurasi dibaca dari file `.env` atau environment variables.

---

### Folder `constants/`

Berisi **konstanta** yang digunakan di seluruh aplikasi, seperti HTTP status codes, error messages, header names, atau nilai-nilai tetap lainnya yang tidak berubah.

---

### Folder `domain/`

Berisi definisi **domain** atau entitas bisnis. Folder ini adalah inti dari clean architecture.

- **`domain/models/`** - Berisi **model/entity** yang merepresentasikan data bisnis (misalnya `User`, `Order`). Model ini digunakan di seluruh lapisan aplikasi.

- **`domain/dto/`** - Berisi **Data Transfer Objects** yang digunakan untuk transfer data antar layer (misalnya request/response structs). Memisahkan struktur data internal dari yang diekspos ke luar.

---

### Folder `repositories/`

Berisi kode untuk **data access layer**. Repository bertanggung jawab untuk komunikasi dengan database atau sumber data lainnya (CRUD operations). Folder ini mengimplementasikan interface yang didefinisikan di domain, memisahkan logika bisnis dari akses data.

---

### Folder `services/`

Berisi **business logic** atau logika bisnis utama aplikasi. Service menggunakan repository untuk mengakses data dan menerapkan aturan bisnis. Layer ini menjadi jembatan antara controller dan repository.

---

### Folder `controllers/`

Berisi **request handler** yang menangani request HTTP masuk. Controller menerima request dari client, memanggil service untuk memproses data, dan mengembalikan response. Folder ini berinteraksi langsung dengan routing.

---

### Folder `middlewares/`

Berisi **HTTP middleware** yang diproses sebelum atau sesudah request sampai ke handler. Contoh penggunaan: autentikasi, autorisasi, logging, CORS, rate limiting, atau validasi token JWT.

---

### Folder `routes/`

Berisi **pendefinisian route/path** untuk endpoint API. Folder ini menghubungkan HTTP method + URL path ke controller/handler yang sesuai.

---

### Folder `clients/`

Berisi kode untuk **HTTP client** atau integrasi dengan **external API**. Digunakan ketika aplikasi perlu berkomunikasi dengan layanan lain (microservice lain, third-party API, dsb).

---

## Alur Request (Clean Architecture)

```
Client Request
      │
      ▼
   Routes  ──▶  Controllers  ──▶  Services  ──▶  Repositories  ──▶  Database
      │              │                 │
      ▼              ▼                 ▼
  (Route)    (Request Handler)  (Business Logic)    (Data Access)
```

1. **Routes** - Menerima request dan mengarahkannya ke controller yang tepat
2. **Controllers** - Memproses request HTTP, memvalidasi input
3. **Services** - Menjalankan logika bisnis
4. **Repositories** - Mengakses dan memanipulasi data di database
5. **Domain** - Definisi model dan DTO yang digunakan di semua layer

---

## Perintah Makefile

| Perintah | Kegunaan |
|----------|----------|
| `make watch-prepare` | Menginstal tool `air` untuk live reload |
| `make watch` | Menjalankan service dengan hot reload (live reload) |
| `make build` | Build/mengkompilasi service menjadi binary |
| `make docker-compose` | Menjalankan service dalam Docker container |
| `make docker-build tag=x.x.x` | Build Docker image dengan tag versi tertentu |
| `make docker-push tag=x.x.x` | Push Docker image ke registry |