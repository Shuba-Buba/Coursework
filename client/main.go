package main

import (
	"client/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func StartCollectTradeInfo(path string, duration time.Duration) {
	SERVER := "stream.binance.com:9443"
	PATH := "ws/btcusdt@trade"
	log.Println("Connecting to:", SERVER, "at", PATH)
	URL := url.URL{Scheme: "wss", Host: SERVER, Path: PATH}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	start := time.Now()

	for duration >= time.Since(start) {
		c, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		defer c.Close()
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Println("Write error:", err)
			return
		}

		tradeInfo := models.TradeInfo{}
		json.Unmarshal(message, &tradeInfo)

		w.Write(tradeInfo.Convert())
	}
	w.Flush()
}

func main() {
	StartCollectTradeInfo("./kek.csv", time.Second*10)
}
