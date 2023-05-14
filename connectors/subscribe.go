package connectors

import (
	"fmt"
	"log"
	"time"
	"trading/binance/http"
	"trading/models"

	"github.com/valyala/fastjson"
)

func SubscribeDepth(port uint, symbol string) <-chan models.Event {
	eventChan := make(chan models.Event, 100)
	go func() {
		defer close(eventChan)

		UDPConn := MakeUDPConnector("224.0.0.1", port)
		parser := fastjson.Parser{}
		buffer := []string{}

		fmt.Print(string(http.GetSnapshot(symbol)))
		// panic("kek")
		for {
			p := make([]byte, 2048)
			size, _, err := UDPConn.ReadFromUDP(p)

			if err != nil {
				log.Panicf("UDPConn error %v", err)
			}

			value, err := parser.Parse(string(p[:size]))
			if err != nil {
				log.Panicf("Error in json parsing: %v", err)
			}
			fmt.Print(value)
			buffer = append(buffer, "aba")
			// if value.GetString<U>

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
