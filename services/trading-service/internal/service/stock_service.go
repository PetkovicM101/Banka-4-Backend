package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/client"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/repository"
)

const maxCallsPerMinute = 55

type StockService struct {
	listingRepo repository.ListingRepository
	stockRepo   repository.StockRepository
	client      *client.StockClient
}

func NewStockService(
	listingRepo repository.ListingRepository,
	stockRepo repository.StockRepository,
	client *client.StockClient,
) *StockService {
	return &StockService{
		listingRepo: listingRepo,
		stockRepo:   stockRepo,
		client:      client,
	}
}

func (s *StockService) SeedStocks(limit int) error {
	symbols, err := s.client.GetSymbols("US")
	if err != nil {
		return fmt.Errorf("failed to fetch symbols: %w", err)
	}

	log.Printf("[seed] fetched %d symbols, seeding up to %d", len(symbols), limit)

	callsThisMinute := 1
	minuteStart := time.Now()
	count := 0

	for _, sym := range symbols {
		if count >= limit {
			break
		}

		if strings.ContainsRune(sym.Symbol, '.') {
			continue
		}

		if callsThisMinute+3 > maxCallsPerMinute {
			elapsed := time.Since(minuteStart)
			if elapsed < time.Minute {
				wait := time.Minute - elapsed
				log.Printf("[seed] rate limit reached, waiting %s...", wait.Round(time.Second))
				time.Sleep(wait)
			}
			callsThisMinute = 0
			minuteStart = time.Now()
		}

		profile, err := s.client.GetProfile(sym.Symbol)
		callsThisMinute++
		if err != nil {
			log.Printf("[seed] skipping %s: profile error: %v", sym.Symbol, err)
			continue
		}
		if profile.Name == "" {
			log.Printf("[seed] skipping %s: empty profile", sym.Symbol)
			continue
		}

		quote, err := s.client.GetQuote(sym.Symbol)
		callsThisMinute++
		if err != nil {
			log.Printf("[seed] skipping %s: quote error: %v", sym.Symbol, err)
			continue
		}
		if quote.CurrentPrice == 0 {
			log.Printf("[seed] skipping %s: no price data", sym.Symbol)
			continue
		}

		financials, err := s.client.GetBasicFinancials(sym.Symbol)
		callsThisMinute++
		if err != nil {
			log.Printf("[seed] skipping %s: financials error: %v", sym.Symbol, err)
			continue
		}

		listing := &model.Listing{
			Ticker:      sym.Symbol,
			Name:        profile.Name,
			ExchangeMIC: profile.Exchange,
			LastRefresh: time.Now(),
			Price:       quote.CurrentPrice,
			Ask:         quote.High,
		}
		if err := s.listingRepo.Upsert(listing); err != nil {
			log.Printf("[seed] skipping %s: listing upsert error: %v", sym.Symbol, err)
			continue
		}

		stock := &model.Stock{
			ListingID:         listing.ListingID,
			OutstandingShares: profile.ShareOutstanding,
			DividendYield:     financials.Metric.DividendYieldIndicatedAnnual,
		}
		if err := s.stockRepo.Upsert(stock); err != nil {
			log.Printf("[seed] skipping %s: stock upsert error: %v", sym.Symbol, err)
			continue
		}

		count++
		log.Printf("[seed] [%d/%d] seeded %s", count, limit, sym.Symbol)
	}

	log.Printf("[seed] done. seeded %d stocks.", count)
	return nil
}
func (s *StockService) Initialize(ctx context.Context) {
	var count int64
	if err := s.listingRepo.Count(&count); err != nil {
		log.Printf("[seed] failed to count listings: %v", err)
		return
	}
	if count == 0 {
		if err := s.SeedStocks(10); err != nil {
			log.Printf("[seed] failed: %v", err)
		}
	}
}

func (s *StockService) StartBackgroundRefresh(ctx context.Context) {
	var count int64
	_ = s.listingRepo.Count(&count)

	if count == 0 {
		go s.StartRefreshLoopNoInitial()
	} else {
		go s.StartRefreshLoop()
	}
}
func (s *StockService) RefreshPrices() error {
	listings, err := s.listingRepo.FindAll()
	if err != nil {
		return fmt.Errorf("failed to load listings: %w", err)
	}

	log.Printf("[refresh] refreshing prices for %d listings", len(listings))

	callsThisMinute := 0
	minuteStart := time.Now()

	for _, listing := range listings {
		if callsThisMinute+1 > maxCallsPerMinute {
			elapsed := time.Since(minuteStart)
			if elapsed < time.Minute {
				wait := time.Minute - elapsed
				log.Printf("[refresh] rate limit reached, waiting %s...", wait.Round(time.Second))
				time.Sleep(wait)
			}
			callsThisMinute = 0
			minuteStart = time.Now()
		}

		quote, err := s.client.GetQuote(listing.Ticker)
		callsThisMinute++
		if err != nil {
			log.Printf("[refresh] skipping %s: %v", listing.Ticker, err)
			continue
		}
		if quote.CurrentPrice == 0 {
			log.Printf("[refresh] skipping %s: no price data", listing.Ticker)
			continue
		}

		if err := s.listingRepo.UpdatePriceAndAsk(&listing, quote.CurrentPrice, quote.High); err != nil {
			log.Printf("[refresh] failed to update %s: %v", listing.Ticker, err)
			continue
		}

		log.Printf("[refresh] updated %s → price=%.4f ask=%.4f", listing.Ticker, quote.CurrentPrice, quote.High)
	}

	log.Printf("[refresh] done")
	return nil
}

func (s *StockService) StartRefreshLoopNoInitial() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.RefreshPrices(); err != nil {
			log.Printf("[refresh] failed: %v", err)
		}
	}
}

func (s *StockService) StartRefreshLoop() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	if err := s.RefreshPrices(); err != nil {
		log.Printf("[refresh] initial run failed: %v", err)
	}

	for range ticker.C {
		if err := s.RefreshPrices(); err != nil {
			log.Printf("[refresh] failed: %v", err)
		}
	}
}
