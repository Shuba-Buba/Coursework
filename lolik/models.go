package lolik

import "time"

type PriceLevel struct {
	Price  float64
	Amount float64
}

func (o *PriceLevel) GetPrice() float64 {
	return o.Price
}

func (o *PriceLevel) GetAmount() float64 {
	return o.Amount
}

type OrderBook struct {
	Bids []PriceLevel
	Asks []PriceLevel
	Time time.Time
}

func (ob *OrderBook) Ask() PriceLevel {
	if len(ob.Asks) > 0 {
		return ob.Asks[0]
	}
	return PriceLevel{}
}

func (ob *OrderBook) Bid() PriceLevel {
	if len(ob.Bids) > 0 {
		return ob.Bids[0]
	}
	return PriceLevel{}
}
