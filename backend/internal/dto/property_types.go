package dto

const PriceMultiplier = 10000

// Property Type - Price
type PTPrice struct {
	Amount   int64  `json:"amount" validate:"gte=0"`
	Currency string `json:"currency" validate:"iso4217"`
}
