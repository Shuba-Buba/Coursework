package exchange

import (
	"trading/common/types"
)

type OrderType string
type OrderSide string
type TimeInForceType string
type RequestOption func(*PlaceOrderParams)

type PlaceOrderParams struct {
	symbol          string
	side            OrderSide
	quantity        float64
	price           float64
	timeInForceType TimeInForceType
	orderType       OrderType
}

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"

	TimeInForceTypeGTC TimeInForceType = "GTC" // Good Till Cancel
	TimeInForceTypeIOC TimeInForceType = "IOC" // Immediate or Cancel
	TimeInForceTypeFOK TimeInForceType = "FOK" // Fill or Kill
	TimeInForceTypeGTX TimeInForceType = "GTX" // Good Till Crossing (Post Only)

	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type PlaceOrderBuilder struct {
	params   *PlaceOrderParams
	exchange Exchange
}

func (t *PlaceOrderBuilder) OrderSide(side OrderSide) *PlaceOrderBuilder {
	t.params.side = side
	return t
}

func (t *PlaceOrderBuilder) Symbol(symbol string) *PlaceOrderBuilder {
	t.params.symbol = symbol
	return t
}

func (t *PlaceOrderBuilder) Quantity(quantity float64) *PlaceOrderBuilder {
	t.params.quantity = quantity
	return t
}

func (t *PlaceOrderBuilder) Price(price float64) *PlaceOrderBuilder {
	t.params.price = price
	return t
}

func (t *PlaceOrderBuilder) TimeInForce(timeInForceType TimeInForceType) *PlaceOrderBuilder {
	t.params.timeInForceType = timeInForceType
	return t
}

func (t *PlaceOrderBuilder) OrderType(orderType OrderType) *PlaceOrderBuilder {
	t.params.orderType = orderType
	return t
}

func (t *PlaceOrderBuilder) Do() (string, error) {
	return t.exchange.PlaceOrder(*t.params)
}

type Exchange interface {
	Init()
	PlaceOrder(PlaceOrderParams) (string, error)
	CancelOrder(symbol string, orderId string) error
	CancelAllOrders(symbol string) error
	GetBalance() error
	GetOrderInfo() OrderInfo
	GetOrderbook() Orderbook
	Update(event types.Event)
	Close()
}
