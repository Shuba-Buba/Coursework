package exchange

import (
	"log"
	"os"
	"trading/common/types"

	"gopkg.in/yaml.v3"
)

type OrderType string
type OrderSide string
type TimeInForceType string

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

type PlaceOrderParams struct {
	symbol          string
	side            OrderSide
	quantity        float64
	price           float64
	timeInForceType TimeInForceType
	orderType       OrderType
}

type NewPlaceOrderRequest struct {
	params   *PlaceOrderParams
	exchange Exchange
}

func (t *NewPlaceOrderRequest) OrderSide(side OrderSide) *NewPlaceOrderRequest {
	t.params.side = side
	return t
}

func (t *NewPlaceOrderRequest) Symbol(symbol string) *NewPlaceOrderRequest {
	t.params.symbol = symbol
	return t
}

func (t *NewPlaceOrderRequest) Quantity(quantity float64) *NewPlaceOrderRequest {
	t.params.quantity = quantity
	return t
}

func (t *NewPlaceOrderRequest) Price(price float64) *NewPlaceOrderRequest {
	t.params.price = price
	return t
}

func (t *NewPlaceOrderRequest) TimeInForce(timeInForceType TimeInForceType) *NewPlaceOrderRequest {
	t.params.timeInForceType = timeInForceType
	return t
}

func (t *NewPlaceOrderRequest) OrderType(orderType OrderType) *NewPlaceOrderRequest {
	t.params.orderType = orderType
	return t
}

func (t *NewPlaceOrderRequest) Do() (string, error) {
	return t.exchange.PlaceOrder(*t.params)
}

func (t *RequestHelper) NewPlaceOrder() *NewPlaceOrderRequest {
	return &NewPlaceOrderRequest{exchange: t.exchange, params: &PlaceOrderParams{}}
}

type CancelOrderParams struct {
	symbol  string
	orderId string
}

type NewCancelOrderRequest struct {
	params   *CancelOrderParams
	exchange Exchange
}

func (t *NewCancelOrderRequest) Symbol(symbol string) *NewCancelOrderRequest {
	t.params.symbol = symbol
	return t
}

func (t *NewCancelOrderRequest) OrderId(orderId string) *NewCancelOrderRequest {
	t.params.orderId = orderId
	return t
}

func (t *NewCancelOrderRequest) Do() error {
	return t.exchange.CancelOrder(*t.params)
}

func (t *RequestHelper) NewCancelOrder() *NewCancelOrderRequest {
	return &NewCancelOrderRequest{exchange: t.exchange, params: &CancelOrderParams{}}
}

type CancelAllOrdersParams struct {
	symbol string
}

type NewCancelAllOrdersRequest struct {
	params   *CancelAllOrdersParams
	exchange Exchange
}

func (t *NewCancelAllOrdersRequest) Symbol(symbol string) *NewCancelAllOrdersRequest {
	t.params.symbol = symbol
	return t
}

func (t *NewCancelAllOrdersRequest) Do() error {
	return t.exchange.CancelAllOrders(*t.params)
}

func (t *RequestHelper) NewCancelAllOrder() *NewCancelAllOrdersRequest {
	return &NewCancelAllOrdersRequest{exchange: t.exchange, params: &CancelAllOrdersParams{}}
}

type RequestHelper struct {
	exchange Exchange
}

type Exchange interface {
	Init()
	PlaceOrder(PlaceOrderParams) (string, error)
	CancelOrder(CancelOrderParams) error
	CancelAllOrders(CancelAllOrdersParams) error
	GetBalance() error
	GetOrderInfo() OrderInfo
	GetOrderbook() Orderbook
	Update(event types.Event)
	Close()
}

type ExchangeClient struct {
	RequestHelper
	Exchange
}

type config struct {
	Binance struct {
		ApiKey    string `yaml:"apiKey"`
		SecretKey string `yaml:"secretKey"`
	}
}

func getConfig() config {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Something went wrong when reading exchange config: %s", err.Error())
	}
	config := config{}
	yaml.Unmarshal(data, &config)
	return config
}

func GetExchangeClient(exchangeType types.ExchangeType) *ExchangeClient {
	config := getConfig()

	var exchange Exchange

	switch exchangeType {
	case types.ExchangeTypeBinance:
		exchange = MakeBinanceExchange(config.Binance.ApiKey, config.Binance.SecretKey)
	}

	return &ExchangeClient{}
}
