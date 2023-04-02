package lolik

type StockMarket struct {
	Name string
}

type Websocket interface {
	SubsribeOnTrades(stock_market StockMarket, callback func(trades []Trades)) error
}
