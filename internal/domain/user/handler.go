package user

import (
	"net/http"
	"strconv"

	"meobeo-talk-api/internal/domain/common"
	"meobeo-talk-api/internal/dto/request"
	"meobeo-talk-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		if err.Error() == "email already exists" {
			common.ErrorResponse(c, http.StatusConflict, common.MessageEmailAlreadyExists, nil)
			return
		}
		common.InternalServerErrorResponse(c)
		return
	}

	common.SuccessResponse(c, common.MessageUserCreated, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			common.NotFoundResponse(c, "User")
			return
		}
		common.InternalServerErrorResponse(c)
		return
	}

	common.SuccessResponse(c, common.MessageUserUpdated, user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		if err.Error() == "user not found" {
			common.NotFoundResponse(c, "User")
			return
		}
		common.InternalServerErrorResponse(c)
		return
	}

	common.SuccessResponse(c, common.MessageUserDeleted, nil)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			common.NotFoundResponse(c, "User")
			return
		}
		common.InternalServerErrorResponse(c)
		return
	}

	common.SuccessResponse(c, common.MessageFetched, user)
}

func (h *Handler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	pagination := &utils.Pagination{
		Page:  page,
		Limit: limit,
	}

	users, err := h.service.GetUsers(pagination)
	if err != nil {
		common.InternalServerErrorResponse(c)
		return
	}

	common.SuccessResponseWithMeta(c, common.MessageUsersFetched, users.Users, users.Meta)
}
