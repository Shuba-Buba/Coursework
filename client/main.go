package main

import (
	"client/models"
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var SERVER string = "stream.binance.com:9443"

func GetUrl(server string, path string) url.URL {
	log.Println("Connecting to:", server, "at", path)
	return url.URL{Scheme: "wss", Host: server, Path: path}
}

func Do(w *csv.Writer, duration time.Duration, URL url.URL, someInfo models.Binance) {
	c, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer c.Close()

	start := time.Now()
	for duration >= time.Since(start) {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Println("Write error:", err)
			return
		}
		someInfo.ParseAndSave(message, w)
	}
}

func StartCollect(path *string, duration time.Duration, URL url.URL, someInfo models.Binance) {
	f, err := os.OpenFile(*path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	Do(w, duration, URL, someInfo)
	w.Flush()
}

func TradeInfo(path string, duration time.Duration) {
	PATH := "ws/btcusdt@trade"
	tmp := models.TradeInfo{}
	StartCollect(&path, duration, GetUrl(SERVER, PATH), tmp)
}

func OrderBookInfo(path string, duration time.Duration) {
	PATH := "ws/btcusdt@depth@100ms"
	tmp := models.OrderBookInfo{}
	StartCollect(&path, duration, GetUrl(SERVER, PATH), tmp)
}

func main() {
	OrderBookInfo("./kek.csv", time.Second*10)
}
