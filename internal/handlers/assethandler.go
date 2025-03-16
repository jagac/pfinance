package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/jagac/pfinance/internal/models"
	"github.com/jagac/pfinance/internal/services"
)

type AssetHandler struct {
	Service          *services.AssetService
	ReturnCalculator *services.ReturnsCalculator
}

func NewAssetHandler(s *services.AssetService, r *services.ReturnsCalculator) *AssetHandler {
	return &AssetHandler{Service: s, ReturnCalculator: r}
}

func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateAsset(r.Context(), &asset); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(asset)
}

func (h *AssetHandler) GetAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := h.Service.Repo.GetAllAssets(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(assets)
}

func (h *AssetHandler) GetReturns(w http.ResponseWriter, r *http.Request) {
	type ReturnsResponse struct {
		StockReturns map[string]float32 `json:"stockReturns"`
		InterestPL   map[string]float32 `json:"interestPl"`
		GoldReturns  map[string]float32 `json:"goldReturns"`
		Error        string             `json:"error,omitempty"`
	}

	var wg sync.WaitGroup

	stockChan := make(chan map[string]float32)
	interestChan := make(chan map[string]float32)
	goldChan := make(chan map[string]float32)
	errorChan := make(chan error, 3)

	wg.Add(3)

	go func() {
		defer wg.Done()
		stockReturns, err := h.ReturnCalculator.StockReturns()
		if err != nil {
			errorChan <- err
			return
		}
		stockChan <- stockReturns
	}()

	go func() {
		defer wg.Done()
		interestPL, err := h.ReturnCalculator.CalculateInterestPL()
		if err != nil {
			errorChan <- err
			return
		}
		interestChan <- interestPL
	}()

	go func() {
		defer wg.Done()
		goldReturns, err := h.ReturnCalculator.GoldReturns()
		if err != nil {
			errorChan <- err
			return
		}
		goldChan <- goldReturns
	}()

	go func() {
		wg.Wait()
		close(stockChan)
		close(interestChan)
		close(goldChan)
		close(errorChan)
	}()

	var response ReturnsResponse
	var hasError bool

	for i := 0; i < 3; i++ {
		select {
		case stockReturns := <-stockChan:
			response.StockReturns = stockReturns
		case interestPL := <-interestChan:
			response.InterestPL = interestPL
		case goldReturns := <-goldChan:
			response.GoldReturns = goldReturns
		case err := <-errorChan:
			hasError = true
			response.Error += err.Error() + " | "
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if hasError {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func (h *AssetHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	asset, err := h.Service.GetAsset(r.Context(), id)
	if err != nil {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(asset)
}
