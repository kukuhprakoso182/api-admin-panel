package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders menyisipkan HTTP security headers ke setiap response
// untuk melindungi dari clickjacking, XSS, sniffing, dan serangan umum lainnya.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cegah browser menebak MIME type (MIME sniffing)
		c.Header("X-Content-Type-Options", "nosniff")

		// Cegah halaman dimuat di dalam iframe (clickjacking)
		c.Header("X-Frame-Options", "DENY")

		// Aktifkan filter XSS bawaan browser lama
		c.Header("X-XSS-Protection", "1; mode=block")

		// Paksa HTTPS untuk semua request berikutnya selama 1 tahun
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Batasi informasi referrer yang dikirim ke domain lain
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Batasi fitur browser yang boleh digunakan
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Hapus header yang bisa bocorkan informasi server
		c.Header("Server", "")
		c.Header("X-Powered-By", "")

		c.Next()
	}
}
