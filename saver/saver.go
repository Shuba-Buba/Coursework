package saver

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"test/contracts"
)

var my_port int = 1234

type Saver struct {
	Port int `json:"port"`
	ch   chan<- contracts.Contract
}

func MakeSaver(sender chan contracts.Contract) *Saver {
	return &Saver{Port: my_port, ch: sender}
}

func Listen(port int, ready chan<- struct{}) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, _ := net.ListenUDP("udp", &addr)

	ready <- struct{}{}

	for {
		p := make([]byte, 2048)
		_, _, err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message %s \n", p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
	}
}

func (this *Saver) Run() {

	wg := sync.WaitGroup{}

	f, err := os.Open("./saver/config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	dec := json.NewDecoder(r)

	addr := net.UDPAddr{
		Port: my_port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}

	ready := make(chan struct{})
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
			Listen(int(port), ready)
			wg.Done()
		}()
		<-ready
		contract.Remote_port = int(port)
		this.ch <- contract
	}
	wg.Wait()
}
