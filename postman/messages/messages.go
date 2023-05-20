package messages

type PostmanRequest struct {
	Type       string `json:"type"`       // one of ["subscription", "heartbeat"]
	Instrument string `json:"instrument"` // binance@futures@BTCUSDT
}

type PostmanResponse struct {
	Instrument string `json:"instrument"`
	Port       uint   `json:"port"`
}
