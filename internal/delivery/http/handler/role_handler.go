package handler

import (
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
	roles, err := h.roleUsecase.GetAll(c.Request.Context())
	if err != nil {
		handleError(c, err, "role")
		return
	}
	response.OK(c, "roles retrieved", roles)
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	role, err := h.roleUsecase.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err, "role")
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
	role, err := h.roleUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err, "role")
		return
	}
	response.Created(c, "role created", role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	var req usecase.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	role, err := h.roleUsecase.Update(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		handleError(c, err, "role")
		return
	}
	response.OK(c, "role updated", role)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	if err := h.roleUsecase.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err, "role")
		return
	}
	response.OK(c, "role deleted", nil)
}