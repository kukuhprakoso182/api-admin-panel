package handler

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"api-admin-panel/internal/domain"
	"api-admin-panel/pkg/response"
)

func handleError(c *gin.Context, err error, resource string) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		response.NotFound(c, resource+" not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		response.Conflict(c, resource+" already exists")
	case errors.Is(err, domain.ErrUnauthorized), errors.Is(err, domain.ErrInvalidCreds):
		response.Unauthorized(c, err.Error())
	case errors.Is(err, domain.ErrForbidden):
		response.Forbidden(c, "access forbidden")
	case errors.Is(err, domain.ErrInvalidReference):
		response.BadRequest(c, "referenced data does not exist", nil)
	case errors.Is(err, domain.ErrInUse):
		response.Conflict(c, resource+" is still in use")
	case errors.Is(err, bcrypt.ErrPasswordTooLong):
		response.BadRequest(c, "password must be at most 72 bytes", nil)
	case errors.Is(err, context.DeadlineExceeded):
		response.GatewayTimeout(c, "request timeout")
	case errors.Is(err, context.Canceled):
		c.AbortWithStatus(499) // klien sudah memutus koneksi
	default:
		// Simpan ke c.Errors agar tercatat oleh HTTPLogger (kolom "Errors").
		_ = c.Error(err)
		response.InternalServerError(c, "internal server error")
	}
}
