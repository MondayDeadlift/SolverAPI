package handler

import (
	"SolverAPI/internal/service"
	"SolverAPI/pkg/codewars"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	CodewarsClient *codewars.Client
	userService    *service.UserService
}

func NewUserHandler(cwClient *codewars.Client, userService *service.UserService) *UserHandler {
	return &UserHandler{
		CodewarsClient: cwClient,
		userService:    userService,
	}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	username := c.Param("username")

	user, err := h.userService.SyncUser(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}
