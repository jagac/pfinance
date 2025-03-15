package jobs

import (
	"context"

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
