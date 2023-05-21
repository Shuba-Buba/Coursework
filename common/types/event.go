package types

import "time"

type EventType int

const (
	OrderBookUpdate EventType = iota
	Trade           EventType = iota
	Snapshot        EventType = iota
)

type Event struct {
	Timestamp time.Time
	Type      EventType
	Data      string
}
