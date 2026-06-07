package models

import (
	"encoding/json"
	"errors"
	"time"
)


type EventType int
type EventSource int

const(
	SERVER EventSource = iota + 1
	CLIENT
)

const(
	MESSAGE EventType = iota + 1
	SUBSCRIBE
	UNSUBSCRIBE
)

type EventDetails struct{
	From json.RawMessage `json:"from"`
	To json.RawMessage `json:"to"`
	Payload json.RawMessage `json:"payload"`
}

type WsEvent struct{
	EventType EventType `json:"event_type"`
	Details EventDetails `json:"details"`
	Source EventSource `json:"src"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *WsEvent) Validate() error {
	if e.EventType < MESSAGE || e.EventType > UNSUBSCRIBE {
		return errors.New("invalid event type")
	}
	if e.Source < SERVER || e.Source > CLIENT {
		return errors.New("invalid event source")
	}
	return nil
}