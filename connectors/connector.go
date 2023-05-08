package connectors

import (
	"fmt"
	"log"
	"net"

	"github.com/gorilla/websocket"
)

type Connector struct {
	Port int
	// Conn      chan string
	// orderBook *Orderbook
}

func (c *Connector) Connect(symbols string) {

	addr := fmt.Sprintf("127.0.0.1:%d", c.Port)
	socket, _, err := websocket.DefaultDialer.Dial(symbols, nil)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer socket.Close()

	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			log.Println("Write error:", err)
			panic(err)
		}
		if err != nil {
			panic("Bad addr")
		}
		conn.Write(message)
	}

}
