package strategy

import "test/exchange"

type Strategy interface {
	OnTick(exchange exchange.Exchange)
	OnFinish()
}
