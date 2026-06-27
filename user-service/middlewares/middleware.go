package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	errConstant "user-service/constants/error"
	services "user-service/services/user"
)

// HandlePanic menangkap panic saat request diproses agar aplikasi tetap mengirim response error standar.
// Middleware ini dipasang secara global agar panic di handler berikutnya tidak membuat server berhenti.
// Return value berupa gin.HandlerFunc karena Gin menjalankan middleware sebagai handler dalam chain request.
func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// defer memastikan recover tetap dijalankan setelah handler berikutnya selesai atau mengalami panic.
		defer func() {
			if r := recover(); r != nil {
				// Panic dicatat ke log lalu response internal server error dikirim ke client.
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstant.ErrInternalServerError.Error(),
				})
				c.Abort()
			}
		}()
		// c.Next menjalankan middleware atau handler berikutnya dalam chain Gin.
		c.Next()
	}
}

// RateLimiter membatasi jumlah request berdasarkan limiter yang diberikan.
// Parameter lmt berisi konfigurasi limit, misalnya jumlah request dan durasi window.
// Jika request melewati limit, middleware mengirim HTTP 429 dan menghentikan chain request.
func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tollbooth mengecek apakah request masih berada dalam batas rate limit.
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errConstant.ErrTooManyRequests.Error(),
			})
			c.Abort()
		}
		// Jika request masih dalam limit, proses diteruskan ke handler berikutnya.
		c.Next()
	}
}

// extractBearerToken mengambil token asli dari format header "Bearer <token>".
// Fungsi ini mengembalikan string kosong jika format Authorization header tidak sesuai.
func extractBearerToken(token string) string {
	// Header dipisah berdasarkan spasi agar bagian "Bearer" dan token bisa dibaca terpisah.
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	return ""
}

// responseUnauthorized mengirim response unauthorized lalu menghentikan request.
// Helper ini dipakai agar semua response 401 memiliki format JSON yang konsisten.
func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	c.Abort()
}

// validateAPIKey memvalidasi API key antar service berdasarkan header dan signature key.
// Header yang dibutuhkan adalah X-Api-Key, X-Request-At, dan X-Service-Name.
// API key dianggap valid jika hash dari serviceName:signatureKey:requestAt sama dengan header X-Api-Key.
func validateAPIKey(c *gin.Context) error {
	// Ambil header yang dikirim oleh service pemanggil.
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	// API key dihitung dari service name, signature key, dan waktu request.
	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	// SHA-256 dipakai untuk membuat signature dari kombinasi data validasi.
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	// Hasil hash dikonversi ke hexadecimal agar bisa dibandingkan dengan header string.
	resultHash := hex.EncodeToString(hash.Sum(nil))

	// Jika API key tidak sama dengan hasil hash, request dianggap tidak valid.
	if apiKey != resultHash {
		return errConstant.ErrUnauthorized
	}
	return nil
}

// validateBearerToken memvalidasi JWT dari Authorization header.
// Parameter token berisi nilai header Authorization lengkap, misalnya "Bearer eyJ...".
// Jika valid, data user dari JWT claims disimpan ke request context.
func validateBearerToken(c *gin.Context, token string) error {
	// Header harus memakai skema Bearer agar token tidak diproses dari format yang salah.
	if !strings.Contains(token, "Bearer") {
		return errConstant.ErrUnauthorized
	}

	// Ambil token JWT tanpa prefix "Bearer".
	tokenString := extractBearerToken(token)
	if tokenString == "" {
		return errConstant.ErrUnauthorized
	}

	// Claims memakai struktur yang sama dengan claims saat token dibuat di service user.
	claims := &services.Claims{}
	// ParseWithClaims memvalidasi token dan membaca data user dari JWT claims.
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Token harus memakai metode signing HMAC.
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errConstant.ErrInvalidToken
		}

		// Secret key harus sama dengan key yang dipakai saat token ditandatangani.
		jwtSecret := []byte(config.Config.JwtSecretKey)
		return jwtSecret, nil
	})

	// Token dianggap gagal jika parsing error atau token tidak valid.
	if err != nil || !tokenJwt.Valid {
		return errConstant.ErrUnauthorized
	}

	// Data user dari claims disimpan ke request context agar bisa dipakai handler berikutnya.
	userLogin := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, claims.User))
	c.Request = userLogin
	// Token juga disimpan di Gin context jika dibutuhkan proses berikutnya.
	c.Set(constants.Token, token)
	return nil
}

// Authenticate memvalidasi Authorization token dan API key sebelum request diteruskan.
// Middleware ini cocok dipasang pada route yang membutuhkan user login dan request antar service yang valid.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		// Authorization header wajib ada untuk endpoint yang membutuhkan autentikasi.
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		// Bearer token divalidasi terlebih dahulu untuk memastikan user terautentikasi.
		err = validateBearerToken(c, token)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		// API key divalidasi untuk memastikan request berasal dari service yang valid.
		err = validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		// Jika semua validasi berhasil, request diteruskan ke handler utama.
		c.Next()
	}
}

/*
Kegunaan file:
File ini dibuat untuk menyimpan middleware HTTP yang dipakai oleh aplikasi.
Middleware ini menangani panic recovery, rate limiting, validasi JWT bearer token,
validasi API key antar service, dan penyimpanan data user login ke request context.
*/
