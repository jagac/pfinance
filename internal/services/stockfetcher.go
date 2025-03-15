package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type StockFetcher struct{}

type StockResponse struct {
	Symbol    string    `json:"symbol"`
	Price     float32   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *StockFetcher) FetchPrice(ticker string) (StockResponse, error) {
	reqUrl := fmt.Sprintf("http://localhost:4000/stock/%s", ticker)

	client := &http.Client{}

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return StockResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return StockResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return StockResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var stockResponse StockResponse
	if err := json.NewDecoder(resp.Body).Decode(&stockResponse); err != nil {
		return StockResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}

	return stockResponse, nil
}
