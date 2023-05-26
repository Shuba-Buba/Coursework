package strategy

import "github.com/Shuba-Buba/Trading-propper-backtest/exchange"

type Strategy interface {
	OnTick(exchange exchange.Exchange)
	OnFinish()
	GetSymbols() []string
}
