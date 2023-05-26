package main

import (
	"time"

	"github.com/Shuba-Buba/Trading-propper-backtest/common/types"
	"github.com/Shuba-Buba/Trading-propper-backtest/exchange"
	"github.com/Shuba-Buba/Trading-propper-backtest/storage"
	"github.com/Shuba-Buba/Trading-propper-backtest/strategy"
	"github.com/Shuba-Buba/Trading-propper-backtest/util"
)

func RunBacktest(fromTime time.Time, toTime time.Time, s strategy.Strategy) {
	var channels []chan types.Event
	for _, v := range s.GetSymbols() {
		channels = append(channels, storage.MakeEventsKeeper(v).GetEvents(fromTime, toTime))
	}
	merged := util.MergeEventChannels(channels)
	exchange := exchange.MakeSimulatedExchange()
	strategy.Run(merged, s, exchange)
}
