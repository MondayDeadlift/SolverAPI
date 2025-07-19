package handler

import (
	"SolverAPI/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type KataHandler struct {
	kataService *service.KataService
}

func NewKataHandler(ks *service.KataService) *KataHandler {
	return &KataHandler{kataService: ks}
}

func (h *KataHandler) GetRandomKata(c echo.Context) error {
	kata, err := h.kataService.GetRandomKata(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to get random kata",
			"details": err.Error(),
		})
	}

	// Возвращаем только необходимые данные
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":   kata.ID,
		"name": kata.Name,
		"url":  kata.URL,
		"tags": kata.Tags,
	})
}
