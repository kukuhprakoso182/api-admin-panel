//go:build integration
// +build integration

package integration_test

// Integration test membutuhkan koneksi database dan server yang berjalan.
// Jalankan dengan: go test ./tests/integration/... -tags=integration -v
import (
	"net/http"
	"os"
	"testing"

	"api-admin-panel/tests/testhelper"

	"github.com/stretchr/testify/assert"
)

var baseURL string

func TestMain(m *testing.M) {
	baseURL = os.Getenv("TEST_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	os.Exit(m.Run())
}

func TestIntegration_Login_Success(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/login", nil)
	req.Header.Set("Content-Type", "application/json")

	payload := map[string]string{
		"email":    "superadmin@example.com",
		"password": "password123",
	}

	req = testhelper.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/login", payload)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("Server tidak tersedia: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestIntegration_Login_WrongPassword(t *testing.T) {
	req := testhelper.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/login", map[string]string{
		"email":    "superadmin@example.com",
		"password": "wrongpassword",
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("Server tidak tersedia: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_ProtectedRoute_WithoutToken(t *testing.T) {
	req := testhelper.NewRequest(http.MethodGet, baseURL+"/api/v1/users", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("Server tidak tersedia: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
