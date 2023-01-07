package models

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
)

type Binance interface {
	Convert() []string
	ParseAndSave(message []byte, w *csv.Writer)
}

type TradeInfo struct {
	EventType        string `json:"e"`
	EventTime        int64  `json:"E"`
	Symbol           string `json:"s"`
	TradeId          int64  `json:"t"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	BuyerOrderId     int64  `json:"b"`
	SellerOrderId    int64  `json:"a"`
	TradeTime        int64  `json:"T"`
	BuyerMarketMaker bool   `json:"m"`
}

func (other *TradeInfo) Convert() []string {
	var record []string
	record = append(
		record, other.EventType, strconv.FormatInt(other.EventTime, 10), other.Symbol,
		strconv.FormatInt(other.TradeId, 10), other.Price, other.Quantity, strconv.FormatInt(other.BuyerOrderId, 10),
		strconv.FormatInt(other.SellerOrderId, 10), strconv.FormatInt(other.TradeTime, 10),
		strconv.FormatBool(other.BuyerMarketMaker))
	return record
}

func (other *TradeInfo) ParseAndSave(message []byte, w *csv.Writer) {
	json.Unmarshal(message, &other)
	w.Write(other.Convert())
}

type PriceLevelAndQuantity struct {
	PriceLevel string
	Quantity   string
}

func parseToString(other *[]PriceLevelAndQuantity) string {
	var result string
	sz := len(*other)
	if sz == 0 {
		return result
	}
	for i := 0; i < sz-1; i++ {
		result += (*other)[i].PriceLevel + "," + (*other)[i].Quantity + ","
	}
	result += (*other)[sz-1].PriceLevel + "," + (*other)[sz-1].Quantity
	return result
}

func (n *PriceLevelAndQuantity) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.PriceLevel, &n.Quantity}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Notification: %d != %d", g, e)
	}
	return nil
}

type OrderBookInfo struct {
	EventType   string                  `json:"e"`
	EventTime   int64                   `json:"E"`
	Symbol      string                  `json:"s"`
	FirstUpdate int64                   `json:"U"`
	FinalUpdate int64                   `json:"u"`
	Bids        []PriceLevelAndQuantity `json:"b"`
	Asks        []PriceLevelAndQuantity `json:"a"`
}

func (other *OrderBookInfo) Convert() []string {
	var record []string
	record = append(record, other.EventType, strconv.FormatInt(other.EventTime, 10),
		other.Symbol, strconv.FormatInt(other.EventTime, 10),
		strconv.FormatInt(other.EventTime, 10), parseToString(&other.Bids),
		parseToString(&other.Asks))
	return record
}

func (other *OrderBookInfo) ParseAndSave(message []byte, w *csv.Writer) {
	json.Unmarshal(message, &other)
	w.Write(other.Convert())
}
