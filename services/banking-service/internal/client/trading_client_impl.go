package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type tradingClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewTradingClient(baseURL string) TradingClient {
	return &tradingClientImpl{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type actuaryProfitResponse struct {
	Profit float64 `json:"profit"`
}

type investmentFundsResponse struct {
	Funds []struct {
		FundName       string  `json:"fund_name"`
		ManagerName    string  `json:"manager_name"`
		BankSharePct   float64 `json:"bank_share_pct"`
		BankShareValue float64 `json:"bank_share_value"`
		Profit         float64 `json:"profit"`
	} `json:"funds"`
}

type actuariesResponse struct {
	Actuaries []struct {
		ID uint `json:"id"`
	} `json:"actuaries"`
}

func (c *tradingClientImpl) GetActuaries(ctx context.Context, authHeader string) ([]Actuary, error) {
	url := fmt.Sprintf("%s/api/actuaries", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+authHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result actuariesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	out := make([]Actuary, 0, len(result.Actuaries))
	for _, a := range result.Actuaries {
		out = append(out, Actuary{ID: a.ID})
	}

	return out, nil
}
func (c *tradingClientImpl) GetActuaryProfit(ctx context.Context, actID uint, authHeader string) (float64, error) {
	url := fmt.Sprintf("%s/api/actuary/%d/assets/profit", c.baseURL, actID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+authHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result actuaryProfitResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Profit, nil
}
func (c *tradingClientImpl) GetInvestmentFunds(ctx context.Context, authHeader string) ([]FundPosition, error) {
	url := fmt.Sprintf("%s/api/investment-funds", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+authHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result investmentFundsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	funds := make([]FundPosition, 0, len(result.Funds))

	for _, f := range result.Funds {
		funds = append(funds, FundPosition{
			FundName:       f.FundName,
			ManagerName:    f.ManagerName,
			BankSharePct:   f.BankSharePct,
			BankShareValue: f.BankShareValue,
			Profit:         f.Profit,
		})
	}

	return funds, nil
}
