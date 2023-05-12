package connectors

import (
	"fmt"
	"log"
	"net"
	"strings"
	"test/contracts"
	"test/models"
	"time"

	"github.com/gorilla/websocket"
)

type Connector struct {
	Port           int
	Symbol         string
	socket_address string
	snapshot       string
	// Conn      chan string
	// orderBook *Orderbook
}

func MakeConnector(other contracts.Contract) (res *Connector) {
	res = &Connector{
		Port:   other.Remote_port,
		Symbol: other.Symbol,
	}

	switch strings.ToLower(other.ExchangeName) {
	case "binance":
		res.socket_address = "wss://stream.binance.com:9443/ws/" + other.Symbol + "@depth@100ms"
		res.snapshot = "Some address!!"
	default:
		panic("Not implemented")
	}
	return
}

func (c *Connector) Connect() {

	addr := fmt.Sprintf("127.0.0.1:%d", c.Port)
	socket, _, err := websocket.DefaultDialer.Dial(c.socket_address, nil)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer socket.Close()

	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			panic("Bad addr")
		}
		conn.Write(message)
	}
}

func (c *Connector) Start() <-chan models.Event {
	eventChan := make(chan models.Event, 100)
	go func() {
		defer close(eventChan)

		addr := net.UDPAddr{
			Port: c.Port,
			IP:   net.ParseIP("127.0.0.1"),
		}
		ser, _ := net.ListenUDP("udp", &addr)

		for {
			p := make([]byte, 2048)
			size, _, err := ser.ReadFromUDP(p)

			if err != nil {
				fmt.Printf("Some error  %v", err)
				continue
			}

			event := models.Event{
				Timestamp: time.Now(),
				EventType: models.OrderBookUpdate,
				Data:      string(p[:size]),
			}

			eventChan <- event
		}
	}()
	return eventChan
}
