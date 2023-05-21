package connectors

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
	"trading/binance/http"
	"trading/common/types"
	"trading/postman/messages"

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

func (this *PostmanConnector) SubscribeDepth(instrument string) <-chan types.Event {
	port, ok := this.Ports[instrument]
	if !ok {
		log.Fatalf("There is no such instrument in PostmanConnector. use `AddInstrument` first.")
	}

	eventChan := make(chan types.Event, 100)
	go func() {
		defer close(eventChan)

		UDPConn := MakeMulticastUDPConnector("224.0.0.1", port)
		parser := fastjson.Parser{}
		p := make([]byte, 65507)
		messageBuffer := []string{}

		var wg sync.WaitGroup
		initMode := true
		prevMessageId := uint64(0)

		for {
			size, _, err := UDPConn.ReadFromUDP(p)
			if err != nil {
				log.Panicf("UDPConn error %v", err)
			}

			// log.Print(string(p[:size]), "size=", size)

			if err != nil {
				log.Panicf("Error in json parsing: %v", err)
			}

			if initMode {
				wg.Wait()
				if initMode {
					// добавляем пришедние обновления ордербука до тех пор пока не получим снепшот
					messageBuffer = append(messageBuffer, string(p[:size]))

					// и после первого пришедшего обновления запрашиваем снепшот
					wg.Add(1)
					go func() {
						snapshotStr := string(http.GetSnapshot(instrument))
						snapshot, err := parser.Parse(snapshotStr)
						if err != nil {
							log.Panic(err)
						}
						e := types.Event{
							Timestamp: time.Now(),
							Type:      types.Snapshot,
							Data:      snapshotStr}
						eventChan <- e

						// time.Sleep()
						lastUpdateId := snapshot.GetUint64("lastUpdateId")
						log.Printf("%s snapshot lastUpdateId = %v", instrument, lastUpdateId)

						log.Printf("len of messageBuffer = %v", len(messageBuffer))
						// парсим все обновления из буффера
						for _, msg := range messageBuffer {

							event, err := parser.Parse(msg)
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
								e := types.Event{
									Timestamp: time.Now(),
									Type:      types.OrderBookUpdate,
									Data:      string(p[:size]),
								}
								eventChan <- e
							}
						}
						initMode = false
						log.Printf("set initMode=%v", initMode)
						wg.Done()
					}()
				}
			}
			if initMode == false {
				eventStr := string(p[:size])
				parsedEvent, err := parser.Parse(eventStr)
				if err != nil {
					log.Fatalf("error occured in json parsing %v", err)
				}
				lastId := parsedEvent.GetUint64("u")
				prevId := parsedEvent.GetUint64("pu")
				if prevMessageId != prevId {
					log.Fatalf("We skipped some messages, prevId = %v, expected %v. Diff = %v", prevId, prevMessageId, prevId-prevMessageId)
				}
				prevMessageId = lastId
				event := types.Event{
					Timestamp: time.Now(),
					Type:      types.OrderBookUpdate,
					Data:      eventStr,
				}
				eventChan <- event
			}
		}
	}()
	return eventChan
}
