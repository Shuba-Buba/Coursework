package exchange

import "trading/common/types"

type Orderbook interface {
	GetAsks(symbol string) []types.Order
	GetBids(symbol string) []types.Order
	Update(event types.Event)
}
