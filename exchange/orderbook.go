package exchange

import "github.com/Shuba-Buba/Trading-propper-backtest/common/types"

type Orderbook interface {
	GetAsks() []types.Order
	GetBids() []types.Order
	Update(event types.Event)
}
