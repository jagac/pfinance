package models

import "time"

type Asset struct {
	ID   int
	Name string `json:"name"`
	Type string `json:"type"`

	Ticker               string    `json:"ticker,omitempty"`
	Price                float32   `json:"price,omitempty"`
	Amount               float32   `json:"amount,omitempty"`
	Currency             string    `json:"currency,omitempty"`
	InterestStart        time.Time `json:"interestStart,omitempty"`
	InterestRate         float32   `json:"interestRate,omitempty"`
	CompoundingFrequency string    `json:"compoundingFrequency,omitempty"`
}
