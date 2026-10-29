package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"api-admin-panel/internal/domain"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/response"
)

func handleError(c *gin.Context, err error, resource string) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		response.NotFound(c, resource+" not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		response.Conflict(c, resource+" already exists")
	case errors.Is(err, usecase.ErrUnauthorized), errors.Is(err, usecase.ErrInvalidCreds):
		response.Unauthorized(c, err.Error())
	case errors.Is(err, usecase.ErrForbidden):
		response.Forbidden(c, "access forbidden")
	case errors.Is(err, context.DeadlineExceeded):
		response.GatewayTimeout(c, "request timeout")
	case errors.Is(err, context.Canceled):
		c.AbortWithStatus(499) // klien sudah memutus koneksi
	default:
		// TODO: log err di sini supaya error tak terduga tetap bisa dilacak
		response.InternalServerError(c, "internal server error")
	}
}