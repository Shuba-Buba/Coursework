package saver

import (
	"sync"
)

func Init(configPath string, postmanPort uint) {

	wg := sync.WaitGroup{}
	saver := MakeSaver(configPath, postmanPort)

	wg.Add(1)

	go func() {
		defer wg.Done()
		saver.Run()
	}()

	wg.Wait()
}
