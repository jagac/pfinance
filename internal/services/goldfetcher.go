package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GoldFetcher struct{}

type GoldResponse struct {
	Price     float64   `json:"xauPrice"`
	Timestamp time.Time `json:"timestamp"`
}

func (g *GoldFetcher) FetchPrice(currency string) (GoldResponse, error) {
	// Format the URL for the specific currency (e.g., USD)
	reqUrl := fmt.Sprintf("https://data-asg.goldprice.org/dbXRates/%s", currency)

	client := &http.Client{}

	// Create the GET request
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return GoldResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return GoldResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GoldResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var goldData struct {
		Items []struct {
			XAUPrice float64 `json:"xauPrice"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&goldData); err != nil {
		return GoldResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}
	pricePerOunce := goldData.Items[0].XAUPrice
	pricePerGram := pricePerOunce / 31.1035

	if len(goldData.Items) == 0 {
		return GoldResponse{}, fmt.Errorf("no price data found for the requested currency")
	}
	goldResponse := GoldResponse{
		Price:     pricePerGram,
		Timestamp: time.Now(),
	}

	return goldResponse, nil
}
