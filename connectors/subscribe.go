package connectors

import (
	"fmt"
	"net"
	"time"
	"trading/models"
)

func Subscribe(port uint) <-chan models.Event {
	eventChan := make(chan models.Event, 100)
	go func() {
		defer close(eventChan)

		addr := net.UDPAddr{
			Port: int(port),
			IP:   net.ParseIP("224.0.0.1"),
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
