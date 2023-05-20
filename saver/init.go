package saver

import (
	"sync"
)

func Init(configPath string, postmanPort uint) {

	wg := sync.WaitGroup{}
	config := ParseConfig(configPath)
	saver := MakeSaver(config, postmanPort)

	wg.Add(1)

	go func() {
		defer wg.Done()
		saver.Run()
	}()

	wg.Wait()
}
