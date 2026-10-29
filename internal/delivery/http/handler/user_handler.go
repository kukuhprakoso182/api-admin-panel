package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-admin-panel/internal/usecase"
	"api-admin-panel/pkg/response"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userUsecase.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		handleError(c, err, "user")
		return
	}
	response.Paginated(c, "users retrieved", users, page, limit, total)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	user, err := h.userUsecase.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err, "user")
		return
	}
	response.OK(c, "user retrieved", user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req usecase.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	user, err := h.userUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err, "user")
		return
	}
	response.Created(c, "user created", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	var req usecase.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	user, err := h.userUsecase.Update(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		handleError(c, err, "user")
		return
	}
	response.OK(c, "user updated", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.userUsecase.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err, "user")
		return
	}
	response.OK(c, "user deleted", nil)
}