package postman

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"trading/connectors"
)

type Postman struct {
	ListenPort    uint
	Connectors    map[string]*connectors.Connector
	FirstFreePort uint
}

func MakePostman(listenPort uint, config PostmanConfig) *Postman {
	return &Postman{
		ListenPort:    listenPort,
		FirstFreePort: config.FirstFreePort,
		Connectors:    make(map[string]*connectors.Connector)}
}

func (this *Postman) Run() {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: int(this.ListenPort),
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	log.Print("Start postman loop")
	for {
		n, remoteaddr, err := conn.ReadFromUDP(p)
		log.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			log.Printf("Error reading from UDP %v", err)
			continue
		}

		var request PostmanRequest
		json.Unmarshal(p[:n], &request)
		log.Print(request)

		switch request.Type {
		case "subscription":
			port := this.GetConnectorPort(request.FullSymbol)
			var response = PostmanResponse{request.FullSymbol, port}
			bytesResponse, _ := json.Marshal(response)
			go conn.WriteToUDP(bytesResponse, remoteaddr)
		case "heartbeat": // раз в 9 минут шлём heatrbeat'ы postman'у, если за 10 минут по Symbol не прилетело ни одного хартбита то выключаем коннектор
			log.Fatal("Not implemented")
		default:
			log.Printf("Unknown request type %s", request.Type)
		}
	}
}

func (this *Postman) GetConnectorPort(fullSymbol string) uint {

	conn, ok := this.Connectors[fullSymbol]
	if ok {
		return conn.Port
	}

	// иначе создаём новый Connector
	port := this.FirstFreePort
	this.FirstFreePort += 1
	splitted := strings.Split(fullSymbol, "@")
	exchange, section, symbol := splitted[0], splitted[1], splitted[2]

	new_connector := connectors.MakeExchangeConnector(exchange, section, symbol, port)
	this.Connectors[fullSymbol] = new_connector

	go new_connector.Connect()

	return new_connector.Port
}
