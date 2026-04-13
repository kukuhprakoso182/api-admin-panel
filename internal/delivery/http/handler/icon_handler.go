package handler

import (
	"errors"

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
	icons, err := h.iconUsecase.GetAll()
	if err != nil {
		response.InternalServerError(c, "failed to retrieve icons")
		return
	}
	response.OK(c, "icons retrieved", icons)
}

func (h *IconHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	icon, err := h.iconUsecase.GetByID(id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "icon not found")
			return
		}
		response.InternalServerError(c, "failed to retrieve icon")
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
	icon, err := h.iconUsecase.Create(&req)
	if err != nil {
		response.InternalServerError(c, "failed to create icon")
		return
	}
	response.Created(c, "icon created", icon)
}

func (h *IconHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req usecase.UpdateIconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	icon, err := h.iconUsecase.Update(id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "icon not found")
			return
		}
		response.InternalServerError(c, "failed to update icon")
		return
	}
	response.OK(c, "icon updated", icon)
}

func (h *IconHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.iconUsecase.Delete(id); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "icon not found")
			return
		}
		response.InternalServerError(c, "failed to delete icon")
		return
	}
	response.OK(c, "icon deleted", nil)
}
