package models

import (
	"encoding/json"
	"errors"
	"strconv"
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
	// To json.RawMessage `json:"to"`
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

func (e *WsEvent) FetchTargetId() (string, error) {
	if e.EventType == MESSAGE {
		var chatmessage ChatMessage
		if err := json.Unmarshal(e.Details.Payload, &chatmessage); err != nil {
			return "", err
		}
		return strconv.Itoa(int(chatmessage.To.Id)), nil
	}

	if e.EventType == SUBSCRIBE || e.EventType == UNSUBSCRIBE {
		var targetId int
		if err := json.Unmarshal(e.Details.Payload, &targetId); err != nil {
			return "", errors.New("invalid payload for subscription: expected integer user id")
		}
		return strconv.Itoa(targetId), nil
	}

	return "", errors.New("unsupported event type for target identification")
}