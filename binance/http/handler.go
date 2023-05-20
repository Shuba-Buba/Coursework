package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type HTTPHandler interface {
	GetSnapshot() []byte
}

func GetSnapshot(instrument string) []byte {
	symbol := strings.Split(instrument, "@")[2]
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/depth?symbol=%s&limit=1000", strings.ToUpper(symbol))

	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}
