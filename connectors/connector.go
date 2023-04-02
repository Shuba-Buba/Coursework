package connectors

import (
	"fmt"
	"log"
	"strings"
	"test/connectors/binance"

	"github.com/gorilla/websocket"
)

type Connector struct {
	// Conn      chan string
	// orderBook *Orderbook
}

func detect(s string) binance.BinanceParser {
	list := strings.Split(s, "@")
	switch list[1] {
	case "depth":
		return &binance.BinanceOrderBook{}
	case "trade":
		return &binance.BinanceTrade{}
	default:
		panic("Not Implemented")
	}

	return nil
}

func (c *Connector) Connect(symbols string) {

	parser := detect(symbols)
	socket, _, err := websocket.DefaultDialer.Dial(symbols, nil)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer socket.Close()

	for {
		_, message, err := socket.ReadMessage()

		if err != nil {
			log.Println("Write error:", err)
			panic("SOSAT")
		}

		parser.Parse(message)
		switch v := parser.(type) {
		default:
			fmt.Println(v)
		}
		fmt.Println(parser.(*binance.BinanceOrderBook))
	}

}
