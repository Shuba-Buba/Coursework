package saver

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"test/connectors"
	"test/contracts"
	"test/storage"
)

var my_port int = 1234

type Saver struct {
	Port int `json:"port"`
	ch   chan<- contracts.Contract
}

func MakeSaver(sender chan contracts.Contract) *Saver {
	return &Saver{Port: my_port, ch: sender}
}

func Listen(current_connector connectors.Connector) {
	events_keeper := storage.MakeEventsKeeper(current_connector.Symbol)
	for event := range current_connector.Start() {
		events_keeper.Save(event)
	}
}

func (this *Saver) Run() {

	wg := sync.WaitGroup{}

	f, err := os.Open("./saver/config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	addr := net.UDPAddr{
		Port: my_port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}

	for {
		contract := contracts.Contract{Port: my_port}
		err := dec.Decode(&contract)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		this.ch <- contract

		p := make([]byte, 100) // какой то большой размер структуры
		var current_connector connectors.Connector

		size, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error %v from %v", err, remoteaddr)
			continue
		}

		err = json.Unmarshal(p[:size], &current_connector)
		if err != nil {
			panic("Bad Connector")
		}

		go func() {
			wg.Add(1)
			Listen(current_connector)
			wg.Done()
		}()
	}
	wg.Wait()
}
