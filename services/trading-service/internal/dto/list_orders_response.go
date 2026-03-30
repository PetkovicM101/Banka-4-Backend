package dto

import "github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"

type OrderSummaryResponse struct {
	OrderID           uint                 `json:"orderId"`
	UserID            uint                 `json:"userId"`
	ListingName       string               `json:"listingName"`
	Quantity          uint                 `json:"quantity"`
	ContractSize      float64              `json:"contractSize"`
	PricePerUnit      *float64             `json:"pricePerUnit"`
	Direction         model.OrderDirection `json:"direction"`
	RemainingPortions uint                 `json:"remainingPortions"`
	Status            model.OrderStatus    `json:"status"`
}

func ToOrderSummaryResponse(o model.Order) OrderSummaryResponse {
	return OrderSummaryResponse{
		OrderID:           o.OrderID,
		UserID:            o.UserID,
		ListingName:       o.Listing.Name,
		Quantity:          o.Quantity,
		ContractSize:      o.ContractSize,
		PricePerUnit:      o.PricePerUnit,
		Direction:         o.Direction,
		RemainingPortions: o.RemainingPortions(),
		Status:            o.Status,
	}
}

func ToOrderSummaryResponseList(orders []model.Order) []OrderSummaryResponse {
	result := make([]OrderSummaryResponse, len(orders))
	for i, o := range orders {
		result[i] = ToOrderSummaryResponse(o)
	}
	return result
}

type ListOrdersResponse struct {
	Data     []OrderSummaryResponse `json:"data"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}
