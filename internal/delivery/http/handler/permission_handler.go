package handler

import (
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
	permissions, err := h.permissionUsecase.GetAll(c.Request.Context())
	if err != nil {
		handleError(c, err, "permission")
		return
	}
	response.OK(c, "permissions retrieved", permissions)
}

func (h *PermissionHandler) GetByID(c *gin.Context) {
	permission, err := h.permissionUsecase.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err, "permission")
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
	permission, err := h.permissionUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err, "permission")
		return
	}
	response.Created(c, "permission created", permission)
}

func (h *PermissionHandler) Update(c *gin.Context) {
	var req usecase.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	permission, err := h.permissionUsecase.Update(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		handleError(c, err, "permission")
		return
	}
	response.OK(c, "permission updated", permission)
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	if err := h.permissionUsecase.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err, "permission")
		return
	}
	response.OK(c, "permission deleted", nil)
}