package types

import (
	"log"
	"strings"
)

type Instrument struct {
	Exchange   ExchangeType // binance
	Section    string       // futures
	Symbol     string       // NEARUSDT
	BaseAsset  string       // NEAR
	QuoteAsset string       // USDT
}

func MakeInstrument(instrument string) Instrument {
	res := Instrument{}
	if len(strings.Split(instrument, "@")) == 3 { // probably binance@futures@btcusdt
		splitted := strings.Split(instrument, "@")
		res.Section = splitted[1]
		res.Symbol = splitted[2]
	} else {
		log.Fatalf("unvalid instrument %s. Use exchange@section@symbol", instrument)
	}
	res.Exchange = ExchangeTypeBinance
	last_cnt := 4 // обычно quote asset это либо USDT либо BUSD, оба имеют 4 символа
	if strings.HasSuffix(res.Symbol, "BTC") {
		last_cnt = 3
	}
	cnt := len(res.Symbol) - last_cnt
	res.BaseAsset = res.Symbol[cnt:]
	res.QuoteAsset = res.Symbol[:cnt]
	return res
}

func (this *Instrument) ToString() string {
	return this.Exchange.ToString() + "@" + this.Section + "@" + this.Symbol
}
