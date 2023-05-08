package exchange

type Exchange interface {
	PlaceOrder()
	CancelOrder()
	CancelAllOrders()
	GetBalance()
	GetOrderInfo()
	GetOrderbook() Orderbook
}
