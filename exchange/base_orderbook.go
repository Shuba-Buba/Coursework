package exchange

import (
	"sort"
	"test/models"
)

type BaseOrderbook struct {
	askPriceToVolume map[float64]float64
	bidPriceToVolume map[float64]float64
	asks             []models.Order
	bids             []models.Order
}

func (t BaseOrderbook) GetAsks() []models.Order {
	return t.asks
}

func (t BaseOrderbook) GetBids() []models.Order {
	return t.bids
}

func (t BaseOrderbook) updateAsks(updates []models.Order) {
	updateMap(t.askPriceToVolume, updates)
	t.asks = updateOrderListFromMap(t.asks, t.askPriceToVolume)
}

func (t BaseOrderbook) updateBids(updates []models.Order) {
	updateMap(t.bidPriceToVolume, updates)
	t.bids = updateOrderListFromMap(t.bids, t.bidPriceToVolume)
}

func (t BaseOrderbook) applySnapshot(asks []models.Order, bids []models.Order) {
	t.askPriceToVolume = make(map[float64]float64)
	t.bidPriceToVolume = make(map[float64]float64)
	updateOrderListFromMap(asks, t.askPriceToVolume)
	updateOrderListFromMap(bids, t.bidPriceToVolume)
}

func updateMap(m map[float64]float64, updates []models.Order) {
	for _, order := range updates {
		if order.Volume == 0 {
			delete(m, order.Price)
		} else {
			m[order.Price] = order.Volume
		}
	}
}

func updateOrderListFromMap(orderList []models.Order, m map[float64]float64) []models.Order {
	var keys []float64
	for k := range m {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	// удалить элементы, но не трогать capacity
	orderList = orderList[:0]

	for _, k := range keys {
		orderList = append(orderList, models.Order{Price: k, Volume: m[k]})
	}

	return orderList
}
