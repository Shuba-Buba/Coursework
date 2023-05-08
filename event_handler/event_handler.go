package event_handler

import (
	"bytes"
	"fmt"
	"net"
	"test/models"
	"time"
)

func StartHandling(port int) <-chan models.Event {
	eventChan := make(chan models.Event, 100)
	go func() {
		addr := net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("127.0.0.1"),
		}
		ser, _ := net.ListenUDP("udp", &addr)

		for {
			p := make([]byte, 2048)
			_, _, err := ser.ReadFromUDP(p)

			if err != nil {
				fmt.Printf("Some error  %v", err)
				continue
			}

			// Отсечение нулевых байт в конце message
			p = bytes.Trim(p, "\x00")

			event := models.Event{
				Timestamp: time.Now(),
				EventType: models.OrderBookUpdate,
				Data:      string(p),
			}

			eventChan <- event
		}
	}()
	return eventChan
}
