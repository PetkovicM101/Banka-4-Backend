package handler

import (
	"net/http"

	"common/pkg/errors"
	"user-service/internal/dto"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service *service.EmployeeService
}

func NewEmployeeHandler(service *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) Register(c *gin.Context) {

	var req dto.CreateEmployeeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	employee, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToEmployeeResponse(employee))
}

func (h *EmployeeHandler) Activate(c *gin.Context) {
	var req dto.ActivateEmployeeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	err := h.service.ActivateAccount(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Password set successfully"})
}
