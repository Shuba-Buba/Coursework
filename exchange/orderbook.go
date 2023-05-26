package exchange

import "trading/common/types"

type Orderbook interface {
	GetAsks() []types.Order
	GetBids() []types.Order
	Update(event types.Event)
}
