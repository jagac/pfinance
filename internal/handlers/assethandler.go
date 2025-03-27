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
	ReturnService    *services.HistoricReturns
}

func NewAssetHandler(s *services.AssetService, r *services.ReturnsCalculator, rs *services.HistoricReturns) *AssetHandler {
	return &AssetHandler{Service: s, ReturnCalculator: r, ReturnService: rs}
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
		Returns map[int]float32 `json:"returns"`
		Error   string          `json:"error,omitempty"`
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	returns := make(map[int]float32)
	errorChan := make(chan error, 3)

	wg.Add(3)

	// Fetch stock returns
	go func() {
		defer wg.Done()
		stockReturns, err := h.ReturnCalculator.StockReturns()
		if err != nil {
			errorChan <- err
			return
		}
		mu.Lock()
		for id, value := range stockReturns {
			returns[id] = value
		}
		mu.Unlock()
	}()

	// Fetch interest returns
	go func() {
		defer wg.Done()
		interestPL, err := h.ReturnCalculator.CalculateInterestPL()
		if err != nil {
			errorChan <- err
			return
		}
		mu.Lock()
		for id, value := range interestPL {
			returns[id] = value
		}
		mu.Unlock()
	}()

	// Fetch gold returns
	go func() {
		defer wg.Done()
		goldReturns, err := h.ReturnCalculator.GoldReturns()
		if err != nil {
			errorChan <- err
			return
		}
		mu.Lock()
		for id, value := range goldReturns {
			returns[id] = value
		}
		mu.Unlock()
	}()

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	var response ReturnsResponse
	response.Returns = returns
	var hasError bool

	for err := range errorChan {
		hasError = true
		response.Error += err.Error() + " | "
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

func (h *AssetHandler) GetMonthlyReturns(w http.ResponseWriter, r *http.Request) {
	returns, err := h.ReturnCalculator.GetMonthlyReturns()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returns)
}
