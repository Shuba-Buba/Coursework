package main

import (
	"time"
	"trading/common/types"
	"trading/exchange"
	"trading/storage"
	"trading/strategy"
	"trading/util"
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
