package models


import (
	"time"
)

type AssetReturn struct {
	ID       int       `json:"id"`
	AssetID  int       `json:"asset_id"`
	Date     time.Time `json:"date"`
	Returns  float64   `json:"returns"`
}
