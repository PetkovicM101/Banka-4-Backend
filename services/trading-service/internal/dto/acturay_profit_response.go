package dto

type ActuaryProfitResponse struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Role      string  `json:"role"`
	ProfitRSD float64 `json:"profitRsd"`
}
