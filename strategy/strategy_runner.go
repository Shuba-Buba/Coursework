package strategy

import (
	"trading/common/types"
	"trading/exchange"
)

func Run(events chan types.Event, strategy Strategy, exchange exchange.Exchange) {

	// ждем снепшот
	for event := range events {
		if event.Type == types.Snapshot {
			exchange.Update(event)
			break
		}
	}

	for event := range events {
		exchange.Update(event)
		strategy.OnTick(exchange)
	}

	strategy.OnFinish()
}
