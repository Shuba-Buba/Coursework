package types

type ExchangeType int64

const (
	Binance ExchangeType = 0
)

func (this *ExchangeType) ToString() string {
	switch *this {
	case Binance:
		return "binance"
	default:
		panic("unknown exchangeType")
	}
	// unreachable
}
