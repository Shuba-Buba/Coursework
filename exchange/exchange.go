package exchange

import "test/models"

type Exchange interface {
	PlaceOrder()
	CancelOrder()
	CancelAllOrders()
	GetBalance()
	GetOrderInfo()
	GetOrderbook() Orderbook
	GetNewTrades()
	Update(event models.Event)
}
