package handler

import (
	"banking-service/internal/dto"
	"banking-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	payment, err := h.service.CreatePayment(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.CreatePaymentResponse{
		ID:     payment.ID,
		Status: string(payment.Status),
	})
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.VerifyPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	payment, err := h.service.VerifyPayment(uint(id), req.Code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.VerifyPaymentResponse{
		ID:     payment.ID,
		Status: string(payment.Status),
	})
}
