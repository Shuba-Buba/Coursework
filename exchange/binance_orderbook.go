package exchange

import (
	"encoding/json"
	"strconv"
	"test/models"
)

type BinanceOrderbook struct {
	BaseOrderbook
}

type binanceOrderbookUpdate struct {
	Bids [][]string `json:"b"`
	Asks [][]string `json:"a"`
}

type binanceOrderbookSnapshot struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func pairsToOrders(pairs [][]string) []models.Order {
	var orders []models.Order
	for i := range pairs {
		price, err := strconv.ParseFloat(pairs[i][0], 64)
		if err != nil {
			panic(err)
		}
		volume, err := strconv.ParseFloat(pairs[i][1], 64)
		orders = append(orders, models.Order{Price: price, Volume: volume})
	}
	return orders
}

func (b *BinanceOrderbook) Update(event models.Event) {
	switch event.Type {
	case models.OrderBookUpdate:
		update := binanceOrderbookUpdate{}
		json.Unmarshal([]byte(event.Data), &update)

		asks := pairsToOrders(update.Asks)
		bids := pairsToOrders(update.Bids)

		b.updateAsks(asks)
		b.updateBids(bids)
	case models.OrderBookSnapshot:
		snapshot := binanceOrderbookSnapshot{}
		json.Unmarshal([]byte(event.Data), &snapshot)

		asks := pairsToOrders(snapshot.Asks)
		bids := pairsToOrders(snapshot.Bids)

		b.applySnapshot(asks, bids)
	default:
		panic("Not implemented")
	}
}
