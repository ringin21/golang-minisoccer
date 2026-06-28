# Tahap builder menyiapkan dependency dan menghasilkan binary user-service.
FROM golang:1.26.3-alpine AS builder

# Install tool sistem yang dibutuhkan selama proses build.
RUN apk update
RUN apk add git openssh tzdata build-base python3 net-tools

# Seluruh proses build dijalankan dari direktori /app.
WORKDIR /app

# Gunakan konfigurasi contoh sebagai konfigurasi build, lalu salin source code.
COPY .env.example .env
COPY . .

# Install Gin CLI, sinkronkan module Go, dan jalankan target build dari Makefile.
RUN go install github.com/buu700/gin@latest
RUN go mod tidy

RUN make build

# Tahap runtime memakai image Alpine yang lebih kecil dari image builder.
FROM alpine:latest

# Siapkan timezone, health-check client, dan direktori kerja aplikasi.
RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    apk --no-cache add curl && \
    mkdir /app

WORKDIR /app

# Dokumentasikan port HTTP yang digunakan user-service.
EXPOSE 8001

# Salin hasil build beserta file aplikasi dari tahap builder.
COPY --from=builder /app /app

# Jalankan user-service sebagai proses utama container.
ENTRYPOINT ["/app/user-service"]
