package saver

import (
	"fmt"
	"sync"
	"trading/contracts"
)

//x.Connect("wss://stream.binance.com:9443/ws/btcusdt@depth@100ms", deliver.Port)

func Start() {
	x := 4 // количество сущностей которые будут поставлять задачи в постман

	ch := make(chan contracts.Contract, x)
	wg := sync.WaitGroup{}
	saver := MakeSaver(ch)

	wg.Add(1)

	go func() {
		defer wg.Done()
		saver.Run()
	}()

	wg.Wait()

}

func Init() {
	Start()
	fmt.Println("OK")
}
