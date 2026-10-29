package handler

import (
	"github.com/gin-gonic/gin"

	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/response"
)

type StatusHandler struct {
	statusUsecase usecase.StatusUsecase
}

func NewStatusHandler(statusUsecase usecase.StatusUsecase) *StatusHandler {
	return &StatusHandler{statusUsecase: statusUsecase}
}

func (h *StatusHandler) GetAll(c *gin.Context) {
	statuses, err := h.statusUsecase.GetAll(c.Request.Context())
	if err != nil {
		handleError(c, err, "status")
		return
	}
	response.OK(c, "statuses retrieved", statuses)
}

func (h *StatusHandler) GetByID(c *gin.Context) {
	status, err := h.statusUsecase.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err, "status")
		return
	}
	response.OK(c, "status retrieved", status)
}

func (h *StatusHandler) Create(c *gin.Context) {
	var req usecase.CreateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	status, err := h.statusUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err, "status")
		return
	}
	response.Created(c, "status created", status)
}

func (h *StatusHandler) Update(c *gin.Context) {
	var req usecase.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	status, err := h.statusUsecase.Update(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		handleError(c, err, "status")
		return
	}
	response.OK(c, "status updated", status)
}

func (h *StatusHandler) Delete(c *gin.Context) {
	if err := h.statusUsecase.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err, "status")
		return
	}
	response.OK(c, "status deleted", nil)
}