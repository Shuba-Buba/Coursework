package postman

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func Init(configPath string, listenPort uint) {
	flag.Parse()

	jsonFile, err := os.Open(configPath)

	if err != nil {
		log.Fatal(err)
	}

	// Парсим конфиг
	byteConfig, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()

	var config PostmanConfig
	json.Unmarshal(byteConfig, &config)
	// use config later ...

	postman := MakePostman(listenPort, config)
	log.Print("Postman created!")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		postman.Run()
	}()

	wg.Wait()
}
