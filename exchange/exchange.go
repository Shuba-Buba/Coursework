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

type PlaceOrderRequestBuilder struct {
	params   *PlaceOrderParams
	exchange Exchange
}

func (t *PlaceOrderRequestBuilder) OrderSide(side OrderSide) *PlaceOrderRequestBuilder {
	t.params.side = side
	return t
}

func (t *PlaceOrderRequestBuilder) Symbol(symbol string) *PlaceOrderRequestBuilder {
	t.params.symbol = symbol
	return t
}

func (t *PlaceOrderRequestBuilder) Quantity(quantity float64) *PlaceOrderRequestBuilder {
	t.params.quantity = quantity
	return t
}

func (t *PlaceOrderRequestBuilder) Price(price float64) *PlaceOrderRequestBuilder {
	t.params.price = price
	return t
}

func (t *PlaceOrderRequestBuilder) TimeInForce(timeInForceType TimeInForceType) *PlaceOrderRequestBuilder {
	t.params.timeInForceType = timeInForceType
	return t
}

func (t *PlaceOrderRequestBuilder) OrderType(orderType OrderType) *PlaceOrderRequestBuilder {
	t.params.orderType = orderType
	return t
}

func (t *PlaceOrderRequestBuilder) Do() (string, error) {
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
