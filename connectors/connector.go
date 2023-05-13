package connectors

import (
	"fmt"
	"log"
	"net"

	"github.com/gorilla/websocket"
)

type Connector struct {
	Port          uint
	Symbol        string
	SocketAddress string
	Snapshot      string
	// Conn      chan string
	// orderBook *Orderbook
}

func MakeExchangeConnector(exchange, section, symbol string, port uint) (conn *Connector) {
	conn = &Connector{
		Port:   port,
		Symbol: symbol,
	}

	if exchange == "binance" && section == "futures" {
		conn.SocketAddress = "wss://fstream.binance.com/ws/" + symbol + "@depth@100ms"
	} else {
		log.Panicf("Connector for %v % v not implemented", exchange, section)
	}
	return
}

func (c *Connector) Connect() {

	addr := fmt.Sprintf("224.0.0.1:%d", c.Port)
	outputConn, err := net.Dial("udp", addr)

	socket, _, err := websocket.DefaultDialer.Dial(c.SocketAddress, nil)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer socket.Close()

	log.Printf("start connector loop for %v at port %v", c.Symbol, c.Port)
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			log.Panic("Couldn't read from websocket")
		}
		outputConn.Write(message)
		log.Printf("Sent message: %s\n", string(message))
	}
}
