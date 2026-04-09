package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	pkgerrors "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
)

type OTCHandler struct {
	service *service.OTCService
}

func NewOTCHandler(service *service.OTCService) *OTCHandler {
	return &OTCHandler{service: service}
}

// PublishAsset godoc
// @Summary Publish assets for OTC trading
// @Description Appends the number of assets the caller makes publicly visible on the OTC portal.
// The amount replaces the current public amount. Must be non-negative and cannot exceed
// owned minus reserved. Accessible by the owning client or actuary only.
// @Tags otc
// @Param id path int true "Asset ownership ID"
// @Param amount path number true "Amount to make public"
// @Success 204
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Security BearerAuth
// @Router /api/client/{clientId}/assets/{id}/publish/{amount} [patch]
func (h *OTCHandler) PublishAsset(c *gin.Context) {
	ownershipID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(pkgerrors.BadRequestErr("invalid asset ownership id"))
		return
	}

	amount, err := strconv.ParseFloat(c.Param("amount"), 64)
	if err != nil {
		c.Error(pkgerrors.BadRequestErr("invalid amount"))
		return
	}

	authCtx := auth.GetAuth(c)
	if authCtx == nil {
		c.Error(pkgerrors.UnauthorizedErr("not authenticated"))
		return
	}

	ownerType := model.OwnerTypeClient
	if authCtx.IdentityType == auth.IdentityEmployee {
		ownerType = model.OwnerTypeActuary
	}

	if svcErr := h.service.PublishAsset(c.Request.Context(), uint(ownershipID), authCtx.IdentityID, ownerType, amount); svcErr != nil {
		c.Error(svcErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPublicOTCAssets godoc
// @Summary List all publicly available OTC assets
// @Description Returns a paginated list of assets that have been marked public by their owners on the OTC portal.
// Each entry includes name, ticker, security type, current price, currency, available amount, last updated timestamp, and owner name.
// @Tags otc
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Security BearerAuth
// @Router /api/otc/public [get]
func (h *OTCHandler) GetPublicOTCAssets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	assets, total, err := h.service.GetPublicOTCAssets(c.Request.Context(), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      assets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
