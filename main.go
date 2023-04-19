package main

import (
	"fmt"
	"sync"
	"test/contracts"
	"test/postman"
	"test/saver"
)

//x.Connect("wss://stream.binance.com:9443/ws/btcusdt@depth@100ms", deliver.Port)

func Start() {
	x := 4 // количетсво сущностей которые будут поставлять задачи в постман

	ch := make(chan contracts.Contract, x)
	wg := sync.WaitGroup{}
	saver := saver.MakeSaver(ch)
	postman := postman.MakePostman(ch)

	wg.Add(2)

	go func() {
		defer wg.Done()
		saver.Run()
	}()

	go func() {
		defer wg.Done()
		postman.Run()
	}()

	wg.Wait()

}

func main() {
	Start()
	fmt.Println("OK")
}
