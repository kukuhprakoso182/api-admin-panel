package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/response"
)

type RoleHandler struct {
	roleUsecase usecase.RoleUsecase
}

func NewRoleHandler(roleUsecase usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{roleUsecase: roleUsecase}
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.roleUsecase.GetAll()
	if err != nil {
		response.InternalServerError(c, "failed to retrieve roles")
		return
	}
	response.OK(c, "roles retrieved", roles)
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	role, err := h.roleUsecase.GetByID(id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "role not found")
			return
		}
		response.InternalServerError(c, "failed to retrieve role")
		return
	}
	response.OK(c, "role retrieved", role)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req usecase.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	role, err := h.roleUsecase.Create(&req)
	if err != nil {
		response.InternalServerError(c, "failed to create role")
		return
	}
	response.Created(c, "role created", role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req usecase.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	role, err := h.roleUsecase.Update(id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "role not found")
			return
		}
		response.InternalServerError(c, "failed to update role")
		return
	}
	response.OK(c, "role updated", role)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.roleUsecase.Delete(id); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "role not found")
			return
		}
		response.InternalServerError(c, "failed to delete role")
		return
	}
	response.OK(c, "role deleted", nil)
}
