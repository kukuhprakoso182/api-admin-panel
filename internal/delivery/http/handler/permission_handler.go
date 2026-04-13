package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/response"
)

type PermissionHandler struct {
	permissionUsecase usecase.PermissionUsecase
}

func NewPermissionHandler(permissionUsecase usecase.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{permissionUsecase: permissionUsecase}
}

func (h *PermissionHandler) GetAll(c *gin.Context) {
	permissions, err := h.permissionUsecase.GetAll()
	if err != nil {
		response.InternalServerError(c, "failed to retrieve permissions")
		return
	}
	response.OK(c, "permissions retrieved", permissions)
}

func (h *PermissionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	permission, err := h.permissionUsecase.GetByID(id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "permission not found")
			return
		}
		response.InternalServerError(c, "failed to retrieve permission")
		return
	}
	response.OK(c, "permission retrieved", permission)
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var req usecase.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	permission, err := h.permissionUsecase.Create(&req)
	if err != nil {
		response.InternalServerError(c, "failed to create permission")
		return
	}
	response.Created(c, "permission created", permission)
}

func (h *PermissionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req usecase.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	permission, err := h.permissionUsecase.Update(id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "permission not found")
			return
		}
		response.InternalServerError(c, "failed to update permission")
		return
	}
	response.OK(c, "permission updated", permission)
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.permissionUsecase.Delete(id); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "permission not found")
			return
		}
		response.InternalServerError(c, "failed to delete permission")
		return
	}
	response.OK(c, "permission deleted", nil)
}
