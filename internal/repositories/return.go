package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/jagac/pfinance/internal/models"
)

type AssetReturnHistoryRepository struct {
	DB *sql.DB
}

func NewAssetReturnHistoryRepository(db *sql.DB) *AssetReturnHistoryRepository {
	return &AssetReturnHistoryRepository{DB: db}
}

func (r *AssetReturnHistoryRepository) InsertAssetReturn(ctx context.Context, assetID int, returns float32, date *time.Time) error {
	var inserted models.AssetReturn
	query := `INSERT INTO asset_returns (asset_id, date, returns) VALUES ($1, COALESCE($2, CURRENT_DATE), $3)`

	err := r.DB.QueryRowContext(ctx, query, assetID, date, returns).Scan(&inserted.ID, &inserted.AssetID, &inserted.Date, &inserted.Returns)
	if err != nil {
		return err
	}
	return nil
}

type MonthlyReturn struct {
	AssetID     int
	AssetName   string
	Year        int
	Month       int
	FirstReturn float64
	LastReturn  float64
}

func (r *AssetReturnHistoryRepository) GetMonthlyReturns(ctx context.Context) ([]MonthlyReturn, error) {
	query := `
	WITH MonthlyReturns AS (
	    SELECT
	        asset_id,
	        DATE_TRUNC('month', date) AS month_start,
	        MIN(date) AS first_date,
	        MAX(date) AS last_date
	    FROM asset_returns
	    GROUP BY asset_id, DATE_TRUNC('month', date)
	), FirstLastReturns AS (
	    SELECT
	        mr.asset_id,
	        EXTRACT(YEAR FROM mr.month_start) AS year,
	        EXTRACT(MONTH FROM mr.month_start) AS month,
	        ar1.returns AS first_return,
	        ar2.returns AS last_return
	    FROM MonthlyReturns mr
	    JOIN asset_returns ar1 ON mr.asset_id = ar1.asset_id AND mr.first_date = ar1.date
	    JOIN asset_returns ar2 ON mr.asset_id = ar2.asset_id AND mr.last_date = ar2.date
	)
	SELECT
	    flr.asset_id,
	    a.name AS asset_name,
	    flr.year,
	    flr.month,
	    flr.first_return,
	    flr.last_return
	FROM FirstLastReturns flr
	JOIN assets a ON flr.asset_id = a.id
	ORDER BY flr.year DESC, flr.month DESC, flr.asset_id;
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var returns []MonthlyReturn
	for rows.Next() {
		var mr MonthlyReturn
		if err := rows.Scan(&mr.AssetID, &mr.AssetName, &mr.Year, &mr.Month, &mr.FirstReturn, &mr.LastReturn); err != nil {
			return nil, err
		}
		returns = append(returns, mr)
	}
	return returns, nil
}
