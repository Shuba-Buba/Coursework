package exchange

import (
	"encoding/json"
	"strconv"

	"github.com/Shuba-Buba/Trading-propper-backtest/common/types"
)

type BinanceOrderbook struct {
	BaseOrderbook
}

func MakeBinanceOrderbook() *BinanceOrderbook {
	return &BinanceOrderbook{*MakeBaseOrderbook()}
}

type binanceOrderbookUpdate struct {
	Bids [][]string `json:"b"`
	Asks [][]string `json:"a"`
}

type binanceSnapshot struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func pairsToOrders(pairs [][]string) []types.Order {
	var orders []types.Order
	for i := range pairs {
		price, err := strconv.ParseFloat(pairs[i][0], 64)
		if err != nil {
			panic(err)
		}
		volume, err := strconv.ParseFloat(pairs[i][1], 64)
		orders = append(orders, types.Order{Price: price, Volume: volume})
	}
	return orders
}

func (b *BinanceOrderbook) Update(event types.Event) {
	switch event.Type {
	case types.EventTypeOrderbookUpdate:
		update := binanceOrderbookUpdate{}
		json.Unmarshal([]byte(event.Data), &update)

		asks := pairsToOrders(update.Asks)
		bids := pairsToOrders(update.Bids)

		b.updateAsks(asks)
		b.updateBids(bids)
	case types.EventTypeSnapshot:
		snapshot := binanceSnapshot{}
		json.Unmarshal([]byte(event.Data), &snapshot)

		asks := pairsToOrders(snapshot.Asks)
		bids := pairsToOrders(snapshot.Bids)

		b.applySnapshot(asks, bids)
	default:
		panic("Not implemented")
	}
}
