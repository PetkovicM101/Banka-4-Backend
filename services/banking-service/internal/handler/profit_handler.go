package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/service"
)

type ProfitHandler struct {
	service *service.ProfitService
}

func NewProfitHandler(service *service.ProfitService) *ProfitHandler {
	return &ProfitHandler{service: service}
}

// GetBankProfit godoc
// @Summary Get bank profit overview
// @Description Returns actuary profits and investment fund positions
// @Tags profit
// @Produce json
// @Success 200 {object} dto.BankProfitResponse
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Security BearerAuth
// @Router /api/bank/profit [get]
func (h *ProfitHandler) GetBankProfit(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	result, err := h.service.GetBankProfit(c.Request.Context(), authHeader)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}
