package postman

type PostmanRequest struct {
	Type       string `json:"type"`   // one of ["subscription", "heartbeat"]
	FullSymbol string `json:"symbol"` // binance@futures@BTCUSDT
}

type PostmanResponse struct {
	Symbol string `json:"symbol"`
	Port   uint   `json:"port"`
}
