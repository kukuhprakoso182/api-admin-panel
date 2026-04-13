package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxRequestSize membatasi ukuran body request yang masuk.
// Melindungi dari serangan DoS dengan mengirim body yang sangat besar.
//
// Contoh penggunaan:
//
//	r.engine.Use(middleware.MaxRequestSize(4 * 1024 * 1024)) // 4MB
func MaxRequestSize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"message": "request body too large",
			})
			c.Abort()
			return
		}

		// Batasi juga saat baca body agar tidak bypass dengan chunked encoding
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
