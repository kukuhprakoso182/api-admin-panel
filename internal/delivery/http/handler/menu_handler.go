package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/response"
)

type MenuHandler struct {
	menuUsecase usecase.MenuUsecase
}

func NewMenuHandler(menuUsecase usecase.MenuUsecase) *MenuHandler {
	return &MenuHandler{menuUsecase: menuUsecase}
}

func (h *MenuHandler) GetAll(c *gin.Context) {
	menus, err := h.menuUsecase.GetAll()
	if err != nil {
		response.InternalServerError(c, "failed to retrieve menus")
		return
	}
	response.OK(c, "menus retrieved", menus)
}

func (h *MenuHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", nil)
		return
	}
	menu, err := h.menuUsecase.GetByID(id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "menu not found")
			return
		}
		response.InternalServerError(c, "failed to retrieve menu")
		return
	}
	response.OK(c, "menu retrieved", menu)
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req usecase.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	menu, err := h.menuUsecase.Create(&req)
	if err != nil {
		response.InternalServerError(c, "failed to create menu")
		return
	}
	response.Created(c, "menu created", menu)
}

func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", nil)
		return
	}
	var req usecase.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	menu, err := h.menuUsecase.Update(id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "menu not found")
			return
		}
		response.InternalServerError(c, "failed to update menu")
		return
	}
	response.OK(c, "menu updated", menu)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", nil)
		return
	}
	if err := h.menuUsecase.Delete(id); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "menu not found")
			return
		}
		response.InternalServerError(c, "failed to delete menu")
		return
	}
	response.OK(c, "menu deleted", nil)
}
