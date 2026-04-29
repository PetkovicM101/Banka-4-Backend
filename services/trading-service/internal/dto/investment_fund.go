package dto

import "time"

type SecurityHoldingResponse struct {
	Ticker            string    `json:"ticker"`
	Price             float64   `json:"price"`
	Change            float64   `json:"change"`
	Volume            uint64    `json:"volume"`
	InitialMarginCost float64   `json:"initialMarginCost"`
	AcquisitionDate   time.Time `json:"acquisitionDate"`
}

type FundPerformanceEntry struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

type FundDetailResponse struct {
	ID                 uint                      `json:"id"`
	Name               string                    `json:"name"`
	Description        string                    `json:"description"`
	Manager            string                    `json:"manager"`
	FundValue          float64                   `json:"fundValue"`
	MinInvestment      float64                   `json:"minInvestment"`
	Profit             float64                   `json:"profit"`
	AccountBalance     float64                   `json:"accountBalance"`
	Holdings           []SecurityHoldingResponse `json:"holdings"`
	PerformanceHistory []FundPerformanceEntry    `json:"performanceHistory"`
}
