package services

import (
	"context"
	"errors"

	"github.com/jagac/pfinance/internal/models"
	"github.com/jagac/pfinance/internal/repositories"
)

type AssetService struct {
	Repo *repositories.AssetRepository
}

func NewAssetService(repo *repositories.AssetRepository) *AssetService {
	return &AssetService{Repo: repo}
}

func (s *AssetService) CreateAsset(ctx context.Context, asset *models.Asset) error {
	if asset.Name == "" {
		return errors.New("asset name is required")
	}
	return s.Repo.AddAsset(ctx, asset)
}

func (s *AssetService) GetAsset(ctx context.Context, id int) (*models.Asset, error) {
	return s.Repo.GetAssetByID(ctx, id)
}
