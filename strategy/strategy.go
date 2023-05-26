package strategy

import "github.com/Shuba-Buba/Trading-propper-backtest/exchange"

type Strategy interface {
	Init(exchange.ExchangeClient)
	OnTick(exchange.ExchangeClient)
	OnFinish(exchange.ExchangeClient)
	GetSymbols() []string
}
