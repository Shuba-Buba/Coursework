package saver

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"test/contracts"
	"test/event_handler"
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

func Listen(symbol string, port int) {
	events_keeper := storage.MakeEventsKeeper(symbol)
	eventChan := event_handler.StartHandling(port)
	for event := range eventChan {
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

		p := make([]byte, 4) // всегда ожидаем инт

		_, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error %v from %v", err, remoteaddr)
			continue
		}

		port := binary.LittleEndian.Uint32(p)
		if err != nil {
			panic(err)
		}
		go func() {
			wg.Add(1)
			Listen(contract.Symbol, int(port))
			wg.Done()
		}()
	}
	wg.Wait()
}
