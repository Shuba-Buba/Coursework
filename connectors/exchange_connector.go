package connectors

import (
	"bytes"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type ExchangeConnector struct {
	Port          uint
	Symbol        string
	SocketAddress string
	Snapshot      string
	// Conn      chan string
	// orderBook *Orderbook
}

func MakeExchangeConnector(exchange, section, symbol string, port uint) (conn *ExchangeConnector) {
	conn = &ExchangeConnector{
		Port:   port,
		Symbol: strings.ToLower(symbol),
	}

	switch exchange {
	case "binance":
		switch section {
		case "futures":
			conn.SocketAddress = "wss://fstream.binance.com/ws/"
		default:
			log.Panicf("ExchangeConnector for %v section %v not implemented", exchange, section)
		}
	default:
		log.Panicf("ExchangeConnector for exchange %v not implemented", exchange)
	}
	return conn
}

func (c *ExchangeConnector) Connect() {

	outputConn := MakeDialUDPConnector("224.0.0.1", c.Port)
	maxMsgLength := 0

	socket, _, err := websocket.DefaultDialer.Dial(c.SocketAddress, nil)
	if err != nil {
		log.Print("Error:", err)
		return
	}
	defer socket.Close()

	startMsg := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{c.Symbol + "@depth@100ms"},
		"id":     1,
	}
	socket.WriteJSON(startMsg)

	log.Printf("start connector loop for %v at port %v", c.Symbol, c.Port)
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			log.Panic("Couldn't read from websocket")
		}
		if bytes.Contains(message, []byte("depthUpdate")) {
			outputConn.Write(message)
			if len(message) > maxMsgLength {
				log.Printf("max message length updated: got len = %v", len(message))
				maxMsgLength = len(message)
			}
		} else {
			log.Printf("Receive message: %s\n", string(message))
		}
	}
}
