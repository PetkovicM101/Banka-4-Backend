package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
)

type ProfitHandler struct {
	service *service.ProfitService
}

func NewProfitHandler(service *service.ProfitService) *ProfitHandler {
	return &ProfitHandler{service: service}
}

// GetActuaryProfits godoc
// @Summary Get actuary profits
// @Description Returns profit for all actuaries (agents/supervisors)
// @Tags profit
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.ActuaryProfitResponse
// @Failure 500 {object} errors.AppError
// @Router /api/profit/actuaries [get]
func (h *ProfitHandler) GetActuaryProfits(c *gin.Context) {
	res, err := h.service.GetActuaryProfits(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetFundPositions godoc
// @Summary Get fund positions
// @Description Returns investment fund positions with bank share and profit
// @Tags profit
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.FundPositionResponse
// @Failure 500 {object} errors.AppError
// @Router /api/profit/funds [get]
func (h *ProfitHandler) GetFundPositions(c *gin.Context) {
	res, err := h.service.GetFundPositions(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
