package common

// TODO: использовать такую хуйню вместо остального говна
type Instrument struct {
	Exchange   string // binance
	Section    string // futures
	Symbol     string // NEARUSDT
	BaseAsset  string // NEAR
	QuoteAsset string // USDT
}
