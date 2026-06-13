package models

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type EventType int
type EventSource int

const (
	SERVER EventSource = iota + 1
	CLIENT
)

const (
	MESSAGE EventType = iota + 1
	SUBSCRIBE
	UNSUBSCRIBE
)

type EventDetails struct {
	From    json.RawMessage `json:"from"`
	Payload json.RawMessage `json:"payload"`
}

type WsEvent struct {
	EventType EventType    `json:"event_type"`
	Details   EventDetails `json:"details"`
	Source    EventSource  `json:"src"`
	Timestamp time.Time    `json:"timestamp"`
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

// NewWsEvent constructs a WsEvent with the given parameters.
// The payload must be a pre-marshaled json.RawMessage.
func NewWsEvent(eType EventType, src EventSource, fromID int64, payload json.RawMessage) WsEvent {
	fromJSON, _ := json.Marshal(map[string]int64{"id": fromID})

	return WsEvent{
		EventType: eType,
		Source:    src,
		Timestamp: time.Now(),
		Details: EventDetails{
			From:    fromJSON,
			Payload: payload,
		},
	}
}

// Marshal converts the WsEvent into a json.RawMessage.
func (e *WsEvent) Marshal() (json.RawMessage, error) {
	return json.Marshal(e)
}
