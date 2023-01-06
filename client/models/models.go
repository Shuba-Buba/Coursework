package models

import "strconv"

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
