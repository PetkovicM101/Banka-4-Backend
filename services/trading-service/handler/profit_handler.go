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


// GetFundPositions godoc
// @Summary Get investment fund positions
// @Description Returns all investment funds with bank share, manager info and profit calculation
// @Tags profit
// @Security BearerAuth
// @Accept json
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
