package exchange

import (
	"sort"

	"github.com/Shuba-Buba/Trading-propper-backtest/common/types"
)

type BaseOrderbook struct {
	askPriceToVolume map[float64]float64
	bidPriceToVolume map[float64]float64
	asks             []types.Order
	bids             []types.Order
}

func MakeBaseOrderbook() *BaseOrderbook {
	return &BaseOrderbook{
		askPriceToVolume: make(map[float64]float64),
		bidPriceToVolume: make(map[float64]float64),
	}
}

func (t BaseOrderbook) GetAsks() []types.Order {
	return t.asks
}

func (t BaseOrderbook) GetBids() []types.Order {
	return t.bids
}

func (t BaseOrderbook) updateAsks(updates []types.Order) {
	updateMap(t.askPriceToVolume, updates)
	t.asks = updateOrderListFromMap(t.asks, t.askPriceToVolume)
}

func (t BaseOrderbook) updateBids(updates []types.Order) {
	updateMap(t.bidPriceToVolume, updates)
	t.bids = updateOrderListFromMap(t.bids, t.bidPriceToVolume)
}

func (t BaseOrderbook) applySnapshot(asks []types.Order, bids []types.Order) {
	t.askPriceToVolume = make(map[float64]float64)
	t.bidPriceToVolume = make(map[float64]float64)
	updateOrderListFromMap(asks, t.askPriceToVolume)
	updateOrderListFromMap(bids, t.bidPriceToVolume)
}

func updateMap(m map[float64]float64, updates []types.Order) {
	for _, order := range updates {
		if order.Volume == 0 {
			delete(m, order.Price)
		} else {
			m[order.Price] = order.Volume
		}
	}
}

func updateOrderListFromMap(orderList []types.Order, m map[float64]float64) []types.Order {
	var keys []float64
	for k := range m {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	// удалить элементы, но не трогать capacity
	orderList = orderList[:0]

	for _, k := range keys {
		orderList = append(orderList, types.Order{Price: k, Volume: m[k]})
	}

	return orderList
}
