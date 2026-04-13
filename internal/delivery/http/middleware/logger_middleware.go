package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"api-admin-panel/pkg/logger"
)

// responseWriter membungkus gin.ResponseWriter agar body response bisa di-capture.
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// HTTPLogger mencatat request dan response dalam format yang rapi dan mudah dibaca.
func HTTPLogger(log *logger.Logger) gin.HandlerFunc {
	skipPaths := map[string]bool{"/health": true}

	return func(c *gin.Context) {
		start     := time.Now()
		path      := c.Request.URL.Path
		rawQuery  := c.Request.URL.RawQuery
		requestID := GetRequestID(c)

		// ── Capture request body ───────────────────────────────────────────────
		var reqBody string
		if !skipPaths[path] && c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			reqBody = sanitizeBody(string(bodyBytes))
		}

		// ── Capture response body ──────────────────────────────────────────────
		rw := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = rw

		c.Next()

		if skipPaths[path] {
			return
		}

		// ── Build log entry ────────────────────────────────────────────────────
		duration   := time.Since(start)
		statusCode := c.Writer.Status()
		fullPath   := path
		if rawQuery != "" {
			fullPath = path + "?" + rawQuery
		}

		message := "request completed"
		switch {
		case statusCode >= 500:
			message = "server error"
		case statusCode >= 400:
			message = "client error"
		}

		errStr := ""
		if len(c.Errors) > 0 {
			errStr = c.Errors.String()
		}

		log.WriteHTTPLog(logger.HTTPLogEntry{
			RequestID:  requestID,
			Method:     c.Request.Method,
			Path:       fullPath,
			UserAgent:  c.Request.UserAgent(),
			ReqBody:    reqBody,
			RespBody:   truncate(rw.body.String(), 4096),
			StatusCode: statusCode,
			LatencyMS:  duration.Milliseconds(),
			Latency:    duration.String(),
			Message:    message,
			Errors:     errStr,
		})
	}
}

// sanitizeBody menyembunyikan field sensitif sebelum ditulis ke log.
func sanitizeBody(body string) string {
	if body == "" {
		return ""
	}
	sensitiveFields := []string{`"password"`, `"token"`, `"secret"`}
	for _, f := range sensitiveFields {
		if strings.Contains(body, f) {
			return "[redacted - contains sensitive fields]"
		}
	}
	return body
}

// truncate memotong string jika melebihi maxLen.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + " ...[truncated]"
}
