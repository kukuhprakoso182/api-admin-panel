package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/response"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Login godoc
// @Summary      Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body usecase.LoginRequest true "Login credentials"
// @Success      200 {object} response.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	tokenPair, err := h.authUsecase.Login(&req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCreds) {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalServerError(c, "login failed")
		return
	}

	response.OK(c, "login successful", tokenPair)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body map[string]string true "Refresh token"
// @Success      200 {object} response.Response
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	tokenPair, err := h.authUsecase.RefreshToken(body.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid or expired refresh token")
		return
	}

	response.OK(c, "token refreshed", tokenPair)
}
