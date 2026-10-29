package handler

import (
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

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", nil)
		return 0, false
	}
	return id, true
}

func (h *MenuHandler) GetAll(c *gin.Context) {
	menus, err := h.menuUsecase.GetAll(c.Request.Context())
	if err != nil {
		handleError(c, err, "menu")
		return
	}
	response.OK(c, "menus retrieved", menus)
}

func (h *MenuHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	menu, err := h.menuUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err, "menu")
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
	menu, err := h.menuUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err, "menu")
		return
	}
	response.Created(c, "menu created", menu)
}

func (h *MenuHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req usecase.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	menu, err := h.menuUsecase.Update(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err, "menu")
		return
	}
	response.OK(c, "menu updated", menu)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.menuUsecase.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err, "menu")
		return
	}
	response.OK(c, "menu deleted", nil)
}