package middleware

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderRequestID  = "X-Request-ID"
const ContextRequestID = "request_id"

// validRequestID membatasi X-Request-ID dari klien agar tidak bisa dipakai
// untuk menyisipkan karakter aneh/sangat panjang ke log.
var validRequestID = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// RequestID menyisipkan X-Request-ID ke setiap request dan response.
// Jika client sudah mengirim header tersebut, nilai itu yang dipakai;
// jika tidak ada atau formatnya tidak valid, di-generate UUID baru.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(HeaderRequestID)
		if !validRequestID.MatchString(rid) {
			rid = uuid.New().String()
		}

		// Simpan ke context agar bisa diambil di handler / middleware lain
		c.Set(ContextRequestID, rid)

		// Kembalikan di response header agar client bisa melacak
		c.Header(HeaderRequestID, rid)

		c.Next()
	}
}

// GetRequestID mengambil request_id dari gin.Context.
func GetRequestID(c *gin.Context) string {
	if rid, ok := c.Get(ContextRequestID); ok {
		if id, ok := rid.(string); ok {
			return id
		}
	}
	return ""
}
