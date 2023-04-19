package connectors

import (
	"log"
	"net"
	"strconv"

	"github.com/gorilla/websocket"
)

type Connector struct {
	Ready         chan struct{}
	Start_working chan struct{}
	Port          int
	// Conn      chan string
	// orderBook *Orderbook
}

func (c *Connector) Connect(symbols string) {

	addr := "127.0.0.1:" + strconv.Itoa(c.Port)
	socket, _, err := websocket.DefaultDialer.Dial(symbols, nil)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer socket.Close()

	c.Ready <- struct{}{}

	<-c.Start_working

	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			log.Println("Write error:", err)
			panic("SOSAT")
		}
		if err != nil {
			panic("Bad addr")
		}
		conn.Write(message)
	}

}
