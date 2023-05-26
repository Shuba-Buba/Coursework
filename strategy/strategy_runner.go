package strategy

import (
	"github.com/Shuba-Buba/Trading-propper-backtest/common/types"
	"github.com/Shuba-Buba/Trading-propper-backtest/exchange"
)

func Run(events chan types.Event, strategy Strategy, exchange exchange.Exchange) {

	// ждем снепшот
	for event := range events {
		if event.Type == types.EventTypeSnapshot {
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
