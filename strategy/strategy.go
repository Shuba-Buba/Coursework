package strategy

import "trading/exchange"

type Strategy interface {
	OnTick(exchange exchange.Exchange)
	OnFinish()
}
