// Package testhelper menyediakan helper umum untuk semua test.
package testhelper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// NewRouter membuat gin.Engine bersih untuk testing.
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

// NewRequest membuat *http.Request dengan JSON body.
func NewRequest(method, path string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewAuthRequest membuat request dengan Bearer token.
func NewAuthRequest(method, path string, body interface{}, token string) *http.Request {
	req := NewRequest(method, path, body)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// Do menjalankan request ke router dan mengembalikan response recorder.
func Do(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ParseBody mem-parse JSON response body ke map.
func ParseBody(w *httptest.ResponseRecorder) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}
