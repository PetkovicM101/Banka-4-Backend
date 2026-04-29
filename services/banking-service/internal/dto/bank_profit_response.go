package dto

type ActuaryProfitResponse struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Role      string  `json:"role"`
	ProfitRSD float64 `json:"profit_rsd"`
}

type FundPositionResponse struct {
	FundName       string  `json:"fund_name"`
	ManagerName    string  `json:"manager_name"`
	BankSharePct   float64 `json:"bank_share_pct"`
	BankShareValue float64 `json:"bank_share_value"`
	ProfitRSD      float64 `json:"profit_rsd"`
}

type BankProfitResponse struct {
	Actuaries []ActuaryProfitResponse `json:"actuaries"`
	Funds     []FundPositionResponse  `json:"funds"`
}
