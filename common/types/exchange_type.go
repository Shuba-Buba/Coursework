package types

type ExchangeType int64

const (
	ExchangeTypeBinance ExchangeType = 0
)

func (this *ExchangeType) ToString() string {
	switch *this {
	case ExchangeTypeBinance:
		return "binance"
	default:
		panic("unknown exchangeType")
	}
	// unreachable
}
