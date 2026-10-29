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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login godoc
// @Summary      Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body usecase.LoginRequest true "Login credentials"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	tokenPair, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err, "user")
		return
	}

	response.OK(c, "login successful", tokenPair)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body RefreshTokenRequest true "Refresh token"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	tokenPair, err := h.authUsecase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			response.Unauthorized(c, "invalid or expired refresh token")
			return
		}
		handleError(c, err, "user")
		return
	}

	response.OK(c, "token refreshed", tokenPair)
}