package strategy

import (
	"test/exchange"
	"test/models"
)

func Run(events chan models.Event, strategy Strategy, exchange exchange.Exchange) {

	// ждем снепшот
	for event := range events {
		if event.Type == models.OrderBookSnapshot {
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
