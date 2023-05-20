package common

// TODO: использовать такое вместо обычных строк
type Instrument struct {
	Exchange   string // binance
	Section    string // futures
	Symbol     string // NEARUSDT
	BaseAsset  string // NEAR
	QuoteAsset string // USDT
}
