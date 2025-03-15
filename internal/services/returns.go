package services

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/jagac/pfinance/internal/repositories"
	"github.com/jagac/pfinance/pkg/cache"
	"github.com/jagac/pfinance/pkg/worker"
)

type ReturnsCalculator struct {
	Repo  *repositories.AssetRepository
	cache *cache.Cache[string, worker.TaskResult]
}

func NewReturnsCalculator(Repo *repositories.AssetRepository,
	cache *cache.Cache[string, worker.TaskResult]) *ReturnsCalculator {
	return &ReturnsCalculator{Repo: Repo, cache: cache}
}

// StockReturns calculates the stock P&L grouped by ticker.
func (r *ReturnsCalculator) StockReturns() (map[string]float32, error) {
	stocks, err := r.Repo.GetAssetsByType(context.Background(), "stock")
	if err != nil {
		return nil, err
	}
	pnlByTicker := make(map[string]float32)

	for _, stock := range stocks {
		currentPrice, ok := r.cache.Get(stock.Ticker)
		if !ok {
			return nil, err
		}
		price, ok := currentPrice.Value.(float32)
		if !ok {
			return nil, err
		}
		pnl := (price - stock.Price) * stock.Amount
		pnlByTicker[stock.Ticker] += pnl
	}
	return pnlByTicker, nil
}

// CalculateInterestPL calculates the interest P&L for assets with interest-bearing properties.
func (r *ReturnsCalculator) CalculateInterestPL() (map[string]float32, error) {
	assets, err := r.Repo.GetAssetsByType(context.Background(), "interest")
	if err != nil {
		return nil, err
	}

	pnlByAsset := make(map[string]float32)

	for _, asset := range assets {
		if asset.InterestRate == 0 || asset.Amount == 0 || asset.InterestStart.IsZero() {
			return nil, errors.New("missing required fields in asset")
		}

		principal := float64(asset.Amount)
		interestRate := float64(asset.InterestRate) / 100
		startDate := asset.InterestStart
		now := time.Now()

		if now.Before(startDate) {
			continue
		}

		monthsElapsed := float64(now.Year()-startDate.Year())*12 + float64(now.Month()-startDate.Month())

		if monthsElapsed <= 0 {
			continue
		}

		var compoundingPeriods int
		switch asset.CompoundingFrequency {
		case "daily":
			compoundingPeriods = 365
		case "quarterly":
			compoundingPeriods = 4
		case "annually":
			compoundingPeriods = 1
		default:
			compoundingPeriods = 12
		}

		yearsElapsed := monthsElapsed / 12
		finalAmount := principal * math.Pow(1+interestRate/float64(compoundingPeriods), float64(compoundingPeriods)*yearsElapsed)
		pl := finalAmount - principal
		pnlByAsset[asset.Name] = float32(pl)
	}

	return pnlByAsset, nil
}

func (r *ReturnsCalculator) GoldReturns() (map[string]float32, error) {
	gold, err := r.Repo.GetAssetsByType(context.Background(), "gold")
	if err != nil {
		return nil, err
	}
	pnlByName := make(map[string]float32)

	for _, g := range gold {
		currentPrice, ok := r.cache.Get("goldPrice")
		if !ok {
			return nil, err
		}
		price, ok := currentPrice.Value.(float32)
		if !ok {
			return nil, err
		}
		pnl := (price - g.Price) * g.Amount
		pnlByName[g.Name] += pnl
	}
	return pnlByName, nil
}
