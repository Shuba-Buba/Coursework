package connectors

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/Shuba-Buba/Trading-propper-backtest/binance/http"
	"github.com/Shuba-Buba/Trading-propper-backtest/models"
	"github.com/Shuba-Buba/Trading-propper-backtest/postman/messages"

	"github.com/valyala/fastjson"
)

type PostmanConnector struct {
	Conn           *net.UDPConn
	PostmanAddress net.UDPAddr
	Ports          map[string]uint
}

func MakePostmanConnector(postmanPort uint) PostmanConnector {
	var res = PostmanConnector{}
	res.Conn = MakeUDPConnector("127.0.0.1", 0)
	res.PostmanAddress = net.UDPAddr{
		Port: int(postmanPort),
		IP:   net.ParseIP("127.0.0.1"),
	}
	res.Ports = make(map[string]uint)
	return res
}

func (this *PostmanConnector) AddInstrument(instrument string) {
	p := make([]byte, 100)

	// отправляем сообщение о подписке постману
	req := messages.PostmanRequest{Type: "subscription", Instrument: instrument}
	b, _ := json.Marshal(req)
	this.Conn.WriteToUDP(b, &this.PostmanAddress)

	// и ждём пока он пришлёт сообщение с портом
	var resp = messages.PostmanResponse{}
	size, _, err := this.Conn.ReadFromUDP(p)
	if err != nil {
		log.Panicf("Error in postman response parsing: %v", err)
	}
	json.Unmarshal(p[:size], &resp)

	if req.Instrument != resp.Instrument {
		log.Fatal("Instruments need to be same")
	}
	this.Ports[instrument] = resp.Port
}

func (this *PostmanConnector) SubscribeDepth(instrument string) <-chan models.Event {
	port, ok := this.Ports[instrument]
	if !ok {
		log.Fatalf("There is no such instrument in PostmanConnector. use `AddInstrument` first.")
	}

	eventChan := make(chan models.Event, 100)
	go func() {
		defer close(eventChan)

		UDPConn := MakeMulticastUDPConnector("224.0.0.1", port)
		parser := fastjson.Parser{}
		messageBuffer := make(chan string, 1000)
		defer close(messageBuffer)

		prevMessageId := uint64(0)

		go func() {
			p := make([]byte, 65507)
			for {
				size, _, err := UDPConn.ReadFromUDP(p)
				if err != nil {
					log.Panicf("UDPConn error %v", err)
				}

				if err != nil {
					log.Panicf("Error in json parsing: %v", err)
				}

				messageBuffer <- string(p[:size])
			}

		}()

		snapshotStr := string(http.GetSnapshot(instrument))
		snapshot, err := parser.Parse(snapshotStr)
		if err != nil {
			log.Panic(err)
		}
		e := models.Event{
			Timestamp: time.Now(),
			Type:      models.Snapshot,
			Data:      snapshotStr,
		}
		eventChan <- e

		lastUpdateId := snapshot.GetUint64("lastUpdateId")
		log.Printf("%s snapshot lastUpdateId = %v", instrument, lastUpdateId)
		for message := range messageBuffer {

			event, err := parser.Parse(message)
			if err != nil {
				log.Panic(err)
			}
			firstId := event.GetUint64("U")
			lastId := event.GetUint64("u")
			prevId := event.GetUint64("pu")
			log.Printf("%s saved event with firstId %v, lastId %v, prevId %v", instrument, firstId, lastId, prevId)

			if prevMessageId != 0 && prevId != prevMessageId {
				log.Panic("We skipped some messages.")
			}
			prevMessageId = lastId

			if lastId >= lastUpdateId {
				e := models.Event{
					Timestamp: time.Now(),
					Type:      models.OrderBookUpdate,
					Data:      message,
				}
				eventChan <- e
			}

		}
	}()
	return eventChan
}
