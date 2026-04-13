package handler

import (
	"errors"

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
	statuses, err := h.statusUsecase.GetAll()
	if err != nil {
		response.InternalServerError(c, "failed to retrieve statuses")
		return
	}
	response.OK(c, "statuses retrieved", statuses)
}

func (h *StatusHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	status, err := h.statusUsecase.GetByID(id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "status not found")
			return
		}
		response.InternalServerError(c, "failed to retrieve status")
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
	status, err := h.statusUsecase.Create(&req)
	if err != nil {
		response.InternalServerError(c, "failed to create status")
		return
	}
	response.Created(c, "status created", status)
}

func (h *StatusHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req usecase.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	status, err := h.statusUsecase.Update(id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "status not found")
			return
		}
		response.InternalServerError(c, "failed to update status")
		return
	}
	response.OK(c, "status updated", status)
}

func (h *StatusHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.statusUsecase.Delete(id); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			response.NotFound(c, "status not found")
			return
		}
		response.InternalServerError(c, "failed to delete status")
		return
	}
	response.OK(c, "status deleted", nil)
}
