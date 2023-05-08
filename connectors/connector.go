package connectors

import (
	"log"
	"time"

	"test/models"

	"github.com/gorilla/websocket"
)

type Connector struct {
	ready  chan chan models.Event
	cancel chan chan models.Event
	// Conn      chan string
	// orderBook *Orderbook
}

func MakeConnector(r, c chan chan models.Event) *Connector {
	return &Connector{ready: r, cancel: c}
}

func (this *Connector) Start() (eventChan <-chan models.Event) {
	eventChan = make(chan models.Event, 100)

	this.ready <- eventChan
	return
}

func (this *Connector) End(rottenEventChan chan models.Event) {
	this.cancel <- rottenEventChan
}

func (c *Connector) Connect(symbols string) {

	socket, _, err := websocket.DefaultDialer.Dial(symbols, nil)

	// addr := fmt.Sprintf("127.0.0.1:%d", c.Port)
	// conn, err := net.Dial("udp", addr)
	// if err != nil {
	// 	log.Println("Error:", err)
	// 	return
	// }
	defer socket.Close()

	subscribes := []chan models.Event{}

	for {
		select {
		case subscribe := <-c.ready:
			subscribes = append(subscribes, subscribe)
		case rottenEventChan := <-c.cancel:
			size := len(subscribes)
			for id, ch := range subscribes {
				if ch == rottenEventChan {
					subscribes[id], subscribes[size-1] = subscribes[size-1], subscribes[id]
					break
				}
			}
			subscribes = subscribes[:size-1]
		default:
			_, message, err := socket.ReadMessage()
			if err != nil {
				log.Println("Write error:", err)
				panic("SOSAT")
			}
			event := models.Event{
				Timestamp: time.Now(),
				EventType: models.OrderBookUpdate,
				Data:      string(message),
			}
			for _, ch := range subscribes {
				ch <- event
			}
			// conn.Write(message)
		}
	}

}
