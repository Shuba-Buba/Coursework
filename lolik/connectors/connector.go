package conectors

import (
	"fmt"
	"log"
	"lolik/conectors/binance"

	"github.com/gorilla/websocket"
)

type Connector struct {
	// Conn      chan string
	// orderBook *Orderbook
}

func (c *Connector) Connect(symbols string) {
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
		parser := binance.Binance{}
		fmt.Println(parser.Parse(message))
	}

}
