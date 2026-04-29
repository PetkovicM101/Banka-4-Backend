package dto

type FundPositionResponse struct {
	FundName       string  `json:"fundName"`
	ManagerName    string  `json:"managerName"`
	BankSharePct   float64 `json:"bankSharePct"`
	BankShareValue float64 `json:"bankShareValue"`
	Profit         float64 `json:"profit"`
}
