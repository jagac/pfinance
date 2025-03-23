package jobs

import (
	"context"

	"github.com/jagac/pfinance/internal/repositories"
	"github.com/jagac/pfinance/internal/services"
	"github.com/jagac/pfinance/pkg/worker"
)

// FetchGoldJob returns a worker job to fetch gold price
func FetchGoldJob(fetcher *services.GoldFetcher) worker.Job {
	return func(c context.Context) (any, error) {
		price, err := fetcher.FetchPrice("EUR")
		if err != nil {
			return nil, err
		}
		return price, nil
	}
}

func FetchStocksJob(assetRepository *repositories.AssetRepository, fetcher *services.StockFetcher) worker.Job {
	return func(c context.Context) (any, error) {
		stocks, err := assetRepository.GetAssetsByType(c, "Stock")
		if err != nil {
			return nil, err
		}

		tickerAndPrice := make(map[string]float32)

		for _, stock := range stocks {
			price, err := fetcher.FetchPrice(stock.Ticker)
			if err != nil {
				return nil, err
			}
			tickerAndPrice[stock.Ticker] = price.Price
		}
		return tickerAndPrice, nil
	}
}

func TotalReturnsJob(returnCalc *services.HistoricReturns) worker.Job {
	return func(c context.Context) (any, error) {
		err := returnCalc.Calc(c)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

}
