package repositories

import (
	"context"
	"database/sql"

	"github.com/jagac/pfinance/internal/models"
)

type AssetRepository struct {
	DB *sql.DB
}

func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{DB: db}
}

func (r *AssetRepository) AddAsset(ctx context.Context, asset *models.Asset) error {
	query := `
		INSERT INTO assets (type, name, ticker, price, amount, currency, interest_rate, compounding_frequency, interest_start)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.DB.ExecContext(ctx, query,
		asset.Type, asset.Name, asset.Ticker, asset.Price, asset.Amount,
		asset.Currency, asset.InterestRate, asset.CompoundingFrequency, asset.InterestStart)

	return err
}

func (r *AssetRepository) GetAssetByID(ctx context.Context, id int) (*models.Asset, error) {
	query := `SELECT * FROM assets WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, query, id)

	var asset models.Asset
	err := row.Scan(&asset.ID, &asset.Type, &asset.Name, &asset.Ticker, &asset.Price, &asset.Amount,
		&asset.Currency, &asset.InterestRate, &asset.CompoundingFrequency, &asset.InterestStart)

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *AssetRepository) GetAllAssets(ctx context.Context) ([]*models.Asset, error) {
	query := `SELECT * FROM assets`
	rows, err := r.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*models.Asset
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(&asset.ID, &asset.Type, &asset.Name, &asset.Ticker, &asset.Price, &asset.Amount,
			&asset.Currency, &asset.InterestRate, &asset.CompoundingFrequency, &asset.InterestStart, &asset.CreatedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *AssetRepository) GetAssetsByType(ctx context.Context, assetType string) ([]*models.Asset, error) {
	query := `SELECT * FROM assets WHERE type = $1`
	rows, err := r.DB.QueryContext(ctx, query, assetType)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*models.Asset
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(&asset.ID, &asset.Type, &asset.Name, &asset.Ticker, &asset.Price, &asset.Amount,
			&asset.Currency, &asset.InterestRate, &asset.CompoundingFrequency, &asset.InterestStart, &asset.CreatedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}
