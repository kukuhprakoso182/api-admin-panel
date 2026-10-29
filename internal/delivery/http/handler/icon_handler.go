package handler

import (
	"github.com/gin-gonic/gin"

	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/response"
)

type IconHandler struct {
	iconUsecase usecase.IconUsecase
}

func NewIconHandler(iconUsecase usecase.IconUsecase) *IconHandler {
	return &IconHandler{iconUsecase: iconUsecase}
}

func (h *IconHandler) GetAll(c *gin.Context) {
	icons, err := h.iconUsecase.GetAll(c.Request.Context())
	if err != nil {
		handleError(c, err, "icon")
		return
	}
	response.OK(c, "icons retrieved", icons)
}

func (h *IconHandler) GetByID(c *gin.Context) {
	icon, err := h.iconUsecase.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err, "icon")
		return
	}
	response.OK(c, "icon retrieved", icon)
}

func (h *IconHandler) Create(c *gin.Context) {
	var req usecase.CreateIconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	icon, err := h.iconUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err, "icon")
		return
	}
	response.Created(c, "icon created", icon)
}

func (h *IconHandler) Update(c *gin.Context) {
	var req usecase.UpdateIconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	icon, err := h.iconUsecase.Update(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		handleError(c, err, "icon")
		return
	}
	response.OK(c, "icon updated", icon)
}

func (h *IconHandler) Delete(c *gin.Context) {
	if err := h.iconUsecase.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err, "icon")
		return
	}
	response.OK(c, "icon deleted", nil)
}