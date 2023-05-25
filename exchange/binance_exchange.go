package exchange

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

type BinanceExchange struct {
	orderTracker *OrderTracker
	client       *futures.Client
	stopUserData chan struct{}
}

func MakeBinanceExchange(apiKey string, secretKey string) *BinanceExchange {
	return &BinanceExchange{
		orderTracker: MakeOrderTracker(),
		client:       binance.NewFuturesClient(apiKey, secretKey),
	}
}

func (t *BinanceExchange) Init() {
	listenKey, err := t.client.NewStartUserStreamService().Do(context.Background())

	if err != nil {
		panic(err)
	}

	// TO DO:
	// Скорее всего здесь нужно запустит горутину,
	// которая будет раз в час будет дергать
	// https://binance-docs.github.io/apidocs/futures/en/#keepalive-user-data-stream-user_stream

	wsHandler := func(event *futures.WsUserDataEvent) {
		if event.Event == futures.UserDataEventTypeOrderTradeUpdate {
			// TO DO: обновить orderInfo в t.orderTracker
		} else if event.Event == futures.UserDataEventTypeAccountConfigUpdate {
			// TO DO: обновить баланс
		} else {
			// TO DO: посмотреть еще бывают типы ивентов
		}
	}

	errHandler := func(err error) {
		fmt.Println("Something wrong with user data channel: ", err)
	}

	_, stopC, err := futures.WsUserDataServe(listenKey, wsHandler, errHandler)
	if err != nil {
		panic(err)
	}

	t.stopUserData = stopC
}

func (t *BinanceExchange) Close() {
	t.stopUserData <- struct{}{}
}

var futuresOrderSide = map[OrderSide]futures.SideType{
	OrderSideBuy:  futures.SideTypeBuy,
	OrderSideSell: futures.SideTypeSell,
}

var futuresTimeInForceType = map[TimeInForceType]futures.TimeInForceType{
	TimeInForceTypeGTC: futures.TimeInForceTypeGTC,
	TimeInForceTypeIOC: futures.TimeInForceTypeIOC,
	TimeInForceTypeGTX: futures.TimeInForceTypeGTX,
	TimeInForceTypeFOK: futures.TimeInForceTypeFOK,
}

var futuresOrderType = map[OrderType]futures.OrderType{
	OrderTypeLimit:  futures.OrderTypeLimit,
	OrderTypeMarket: futures.OrderTypeMarket,
}

func (t *BinanceExchange) PlaceOrder(params PlaceOrderParams) (string, error) {
	order, err := t.client.NewCreateOrderService().
		Symbol(params.symbol).
		Side(futuresOrderSide[params.side]).
		Type(futuresOrderType[params.orderType]).
		TimeInForce(futuresTimeInForceType[params.timeInForceType]).
		Quantity(strconv.FormatFloat(params.quantity, 'f', 10, 64)).
		Price(strconv.FormatFloat(params.price, 'f', 10, 64)).
		Do(context.Background())

	if err != nil {
		return "", err
	}

	return order.ClientOrderID, nil
}

func (t *BinanceExchange) CancelOrder(params CancelOrderParams) error {
	// Игнорируется ответ, хотя наверное в нем есть что-то полезное, что можно вернуть
	_, err := t.client.NewCancelOrderService().
		Symbol(params.symbol).
		OrigClientOrderID(params.orderId). // ?? вот здесь origClientOrderId это clientOrderId, который был получен при PlaceOrder??
		Do(context.Background())
	return err
}

func (t *BinanceExchange) CancelAllOrders(params CancelAllOrdersParams) error {
	err := t.client.NewCancelAllOpenOrdersService().Symbol(params.symbol).Do(context.Background())
	return err
}
