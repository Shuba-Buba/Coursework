package connectors

import (
	"encoding/json"
	"log"
	"net"
	"time"
	"trading/models"
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

func (this *PostmanConnector) SubscribeDepth(instrument string) <-chan models.Event {
	port, ok := this.Ports[instrument]
	if !ok {
		log.Fatalf("There is")
	}

	eventChan := make(chan models.Event, 100)
	go func() {
		defer close(eventChan)

		UDPConn := MakeMulticastUDPConnector("224.0.0.1", port)
		parser := fastjson.Parser{}
		p := make([]byte, 2048)
		buffer := []string{}

		// fmt.Print(string(http.GetSnapshot(instrument)))
		log.Printf("port = %v", port)

		for {
			size, _, err := UDPConn.ReadFromUDP(p)
			if err != nil {
				log.Panicf("UDPConn error %v", err)
			}

			value, err := parser.Parse(string(p[:size]))
			if err != nil {
				log.Panicf("Error in json parsing: %v", err)
			}

			value.Array()
			// log.Print(value)
			log.Print(string(p[:size]))
			buffer = append(buffer, "aba")
			// if value.GetString<U>

			event := models.Event{
				Timestamp: time.Now(),
				// EventType: models.OrderBookUpdate,
				Data: string(p[:size]),
			}

			eventChan <- event
		}
	}()
	return eventChan
}
