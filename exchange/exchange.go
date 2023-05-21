package exchange

import "trading/common/types"

type Exchange interface {
	PlaceOrder()
	CancelOrder()
	CancelAllOrders()
	GetBalance()
	GetOrderInfo()
	GetOrderbook() Orderbook
	GetNewTrades()
	Update(event types.Event)
}
