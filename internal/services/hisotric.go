package services

import (
	"context"
	"math"
	"time"

	"github.com/jagac/pfinance/internal/repositories"
)

type HistoricReturns struct {
	assetRepo    *repositories.AssetRepository
	historicRepo *repositories.AssetReturnHistoryRepository
	goldFetcher  *GoldFetcher
	stockFetcher *StockFetcher
}

func NewHistoricReturns(assetRepo *repositories.AssetRepository,
	historicRepo *repositories.AssetReturnHistoryRepository,
	goldFetcher *GoldFetcher,
	stockFetcher *StockFetcher) *HistoricReturns {
	return &HistoricReturns{assetRepo: assetRepo, historicRepo: historicRepo, goldFetcher: goldFetcher, stockFetcher: stockFetcher}
}

func (r *HistoricReturns) Calc(ctx context.Context) error {
	assets, err := r.assetRepo.GetAllAssets(ctx)
	if err != nil {
		return err
	}

	for _, asset := range assets {
		if asset.Type == "Stock" {

			stockPrice, err := r.stockFetcher.FetchPrice(asset.Ticker)
			if err != nil {
				return err
			}

			currentPrice32 := float32(stockPrice.Price)
			pnl := (currentPrice32 - asset.Price) * asset.Amount
			err = r.historicRepo.InsertAssetReturn(ctx, asset.ID, pnl, &stockPrice.Timestamp)
			if err != nil {
				return err
			}
		}

		if asset.Type == "Gold" {

			currentGoldPrice, err := r.stockFetcher.FetchPrice(asset.Ticker)
			if err != nil {
				return err
			}
			currentPrice32 := float32(currentGoldPrice.Price)
			pnl := (currentPrice32 - asset.Price) * asset.Amount
			err = r.historicRepo.InsertAssetReturn(ctx, asset.ID, pnl, &currentGoldPrice.Timestamp)
			if err != nil {
				return err
			}
		}

		if asset.Type == "Savings" {
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
			pl32 := float32(pl)

			err = r.historicRepo.InsertAssetReturn(ctx, asset.ID, pl32, &now)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
