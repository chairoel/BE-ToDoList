package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"todolist/internal/domain"
	"todolist/internal/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) Create(c *gin.Context) {
	var input domain.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		badRequest(c, err.Error())
		return
	}

	user, err := h.userUsecase.Create(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	successResponse(c, http.StatusCreated, "user created successfully", user)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userUsecase.GetAll(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	successResponse(c, http.StatusOK, "users fetched successfully", users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "invalid user id")
		return
	}

	user, err := h.userUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	successResponse(c, http.StatusOK, "user fetched successfully", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "invalid user id")
		return
	}

	var input domain.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		badRequest(c, err.Error())
		return
	}

	user, err := h.userUsecase.Update(c.Request.Context(), id, input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	successResponse(c, http.StatusOK, "user updated successfully", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "invalid user id")
		return
	}

	if err := h.userUsecase.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	successResponse(c, http.StatusOK, "user deleted successfully", nil)
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		notFound(c, err.Error())
	case errors.Is(err, domain.ErrEmailAlreadyUsed):
		conflict(c, err.Error())
	case errors.Is(err, domain.ErrInvalidUserInput):
		badRequest(c, err.Error())
	default:
		internalError(c, "internal server error")
	}
}
