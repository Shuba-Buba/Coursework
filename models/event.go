package models

import "time"

type EventType int

const (
	OrderBookUpdate EventType = iota
	Trade           EventType = iota
)

type Event struct {
	Timestamp time.Time
	EventType EventType
	Data      string
}
