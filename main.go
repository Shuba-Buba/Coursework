package main

import (
	"test/connectors"
)

func main() {
	x := connectors.Connector{}
	x.Connect("wss://stream.binance.com:9443/ws/btcusdt@depth@100ms")
}
