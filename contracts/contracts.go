package contracts

type Contract struct {
	Symbol      string `json:"symbol"`
	MarketName  string `json:"stock_market_name"` // не придумал ничего умнее
	Port        int
	Remote_port int
}
