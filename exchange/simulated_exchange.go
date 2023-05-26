package exchange

import (
	"strconv"
	"trading/common/types"
)

// Init()
// PlaceOrder(PlaceOrderParams) (string, error)
// CancelOrder(CancelOrderParams) error
// CancelAllOrders(CancelAllOrdersParams) error
// GetBalance(symbol string) float64
// GetOrderInfo(orderId string) OrderInfo
// GetOrderbook(symbol string) Orderbook
// Update(event types.Event)
// Close()

type OrderWithId struct {
	id    string
	order types.Order
}

type SimulatedExchange struct {
	orderTracker OrderTracker
	orderbooks   map[string]Orderbook
	orders       map[string][]OrderWithId
	nextOrderId  int64
}

func (t *SimulatedExchange) Init() {
	// do nothing
}

func (t *SimulatedExchange) PlaceOrder(params PlaceOrderParams) (string, error) {
	orderId := strconv.FormatInt(t.nextOrderId, 10)
	t.nextOrderId++

	orderWithId := OrderWithId{
		id:    orderId,
		order: types.Order{Price: params.price, Volume: params.quantity},
	}

	t.orders[params.symbol] = append(t.orders[params.symbol], orderWithId)

	return orderId, nil
}

func (t *SimulatedExchange) CancelOrder(params CancelOrderParams) error {
	orders := t.orders[params.symbol]
	for i, v := range orders {
		if v.id == params.orderId {
			orders[i] = orders[len(orders)-1]
			t.orders[params.symbol] = orders[:len(orders)-1]
			break
		}
	}
	return nil
}

func (t *SimulatedExchange) CancelAllOrders(params CancelAllOrdersParams) error {
	delete(t.orders, params.symbol)
	return nil
}

func (t *SimulatedExchange) GetOrderbook(symbol string) Orderbook {
	return t.orderbooks[symbol]
}

func (t *SimulatedExchange) GetBalance(symbol string) float64 {
	panic("Not implemented")
}

func (t *SimulatedExchange) GetOrderInfo(orderId string) OrderInfo {
	panic("Not implemented")
}

func (t *SimulatedExchange) Update(event types.Event) {
	panic("Not implemented")
}

func (t *SimulatedExchange) Close() {
	// do nothing
}

func MakeSimulatedExchange() *SimulatedExchange {
	return &SimulatedExchange{
		orderTracker: OrderTracker{},
		orderbooks:   map[string]Orderbook{},
	}
}
