package exchange

import "test/models"

type Orderbook interface {
	GetAsks(symbol string) []models.Order
	GetBids(symbol string) []models.Order
	Update(event models.Event)
}

