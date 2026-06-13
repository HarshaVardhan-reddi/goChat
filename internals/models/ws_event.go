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

// Message is the interface for all payloads (Chat, Status, etc.)
type Message interface {
	GetFromID() int64
	GetToID() int64
}

// StatusMessage implements the Message interface
type StatusMessage struct {
	UserID   int64     `json:"user_id"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

func (s StatusMessage) GetFromID() int64 { return s.UserID }
func (s StatusMessage) GetToID() int64   { return 0 } // Subscriptions are often broadcast or logic-specific

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

// GetMessage unmarshals the raw payload into the correct concrete type
func (e *WsEvent) GetMessage() (Message, error) {
	switch e.EventType {
	case MESSAGE:
		var m ChatMessage
		err := json.Unmarshal(e.Details.Payload, &m)
		return m, err
	case SUBSCRIBE, UNSUBSCRIBE:
		// Try parsing as integer first (client style)
		var targetId int64
		if err := json.Unmarshal(e.Details.Payload, &targetId); err == nil {
			return StatusMessage{UserID: targetId}, nil
		}
		// Fallback to StatusMessage struct (server style)
		var m StatusMessage
		err := json.Unmarshal(e.Details.Payload, &m)
		return m, err
	default:
		return nil, errors.New("unknown event type")
	}
}

func (e *WsEvent) FetchTargetId() (string, error) {
	msg, err := e.GetMessage()
	if err != nil {
		return "", err
	}

	if e.EventType == MESSAGE {
		return strconv.FormatInt(msg.GetToID(), 10), nil
	}

	if e.EventType == SUBSCRIBE || e.EventType == UNSUBSCRIBE {
		return strconv.FormatInt(msg.GetFromID(), 10), nil
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
