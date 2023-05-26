# Trading-propper-backtest

Module for backtesting trading models.

# How To Use

1. Write your own implementation of strategy.Strategy intreface.
2. Call RunBacktest.

# Example

```
type MyStrategy struct{}

const (
	BTCUSDT_SYMBOL string = "BTCUSDT"
)

func (t *MyStrategy) Init(exchange exchange.ExchangeClient) {
	// do something before start
}

func (t *MyStrategy) GetSymbols() []string {
	return []string{BTCUSDT_SYMBOL}
}

func (t *MyStrategy) OnTick(client exchange.ExchangeClient) {
	if client.GetOrderbook(BTCUSDT_SYMBOL).GetAsks()[0].Price < 10.0 {
		orderId, err := client.NewPlaceOrder().Symbol(BTCUSDT_SYMBOL).OrderSide(exchange.OrderSideBuy).
			Price(42.0).Quantity(13.0).Do()
		err = client.NewCancelOrder().Symbol(BTCUSDT_SYMBOL).
			OrderId(orderId).Do()
	}

	if client.GetBalance(BTCUSDT_SYMBOL) < 20.0 {
		err := client.NewCancelAllOrder.Symbol(BTCUSDT_SYMBOL).Do()
	}
}

func (t *MyStrategy) OnFinish(exchange exchange.ExchangeClient) {
	exchange.NewCancelAllOrder().Symbol(BTCUSDT_SYMBOL).Do()
}

func main() {
	fromTime := time.Now().Add(-time.Hour * 24)
	toTime := time.Now()
	RunBacktest(fromTime, toTime, &MyStrategy{})
}

```