package contracts

type Contract struct {
	Symbol       string `json:"symbol"`
	ExchangeName string `json:"stock_market_name"`
	Port         int
	Remote_port  int
}
