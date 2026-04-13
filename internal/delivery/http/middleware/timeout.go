package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout membatalkan request yang melebihi durasi yang ditentukan.
// Melindungi dari request yang tergantung terlalu lama (slowloris, query berat).
//
// Contoh penggunaan:
//
//	r.engine.Use(middleware.Timeout(30 * time.Second))
func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()

		// Ganti context request dengan yang ber-timeout
		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{}, 1)

		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			// Request selesai sebelum timeout
		case <-ctx.Done():
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"success": false,
				"message": "request timeout",
			})
			c.Abort()
		}
	}
}
