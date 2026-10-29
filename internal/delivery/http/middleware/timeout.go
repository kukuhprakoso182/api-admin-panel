package middleware

import (
	"context"
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
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
