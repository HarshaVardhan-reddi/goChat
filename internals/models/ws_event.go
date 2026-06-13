package models

import (
	"chatonetoone/internals/models/events"
	"encoding/json"
	"errors"
	"time"
)

type EventType int

const (
	MESSAGE EventType = iota + 1
	SUBSCRIBE
	UNSUBSCRIBE
	STATUS_UPDATE
)

type EventSource int

const (
	SERVER EventSource = iota + 1
	CLIENT
)

// --- WsEvent Structure ---

type EventDetails struct {
	Type    EventType      `json:"type"`
	Message events.Message `json:"message"`
}

type WsEvent struct {
	From      events.UserRef `json:"from"`
	To        events.UserRef `json:"to"`
	Source    EventSource    `json:"src"`
	Details   EventDetails   `json:"details"`
	Timestamp time.Time      `json:"timestamp"`
}

// Custom Unmarshaler for EventDetails to handle the Message interface polymorphism
func (ed *EventDetails) UnmarshalJSON(data []byte) error {
	type Alias EventDetails
	aux := &struct {
		RawMessage json.RawMessage `json:"message"`
		*Alias
	}{
		Alias: (*Alias)(ed),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch ed.Type {
	case MESSAGE:
		var m events.ChatMessage
		if err := json.Unmarshal(aux.RawMessage, &m); err != nil {
			return err
		}
		ed.Message = m
	case SUBSCRIBE:
		var m events.SubscribeMessage
		if err := json.Unmarshal(aux.RawMessage, &m); err != nil {
			// Fallback: Client might just send the ID as an integer
			var targetId int64
			if err2 := json.Unmarshal(aux.RawMessage, &targetId); err2 == nil {
				ed.Message = events.SubscribeMessage{TargetID: targetId}
				return nil
			}
			return err
		}
		ed.Message = m
	case UNSUBSCRIBE:
		var m events.UnsubscribeMessage
		if err := json.Unmarshal(aux.RawMessage, &m); err != nil {
			var targetId int64
			if err2 := json.Unmarshal(aux.RawMessage, &targetId); err2 == nil {
				ed.Message = events.UnsubscribeMessage{TargetID: targetId}
				return nil
			}
			return err
		}
		ed.Message = m
	case STATUS_UPDATE:
		var m events.StatusMessage
		if err := json.Unmarshal(aux.RawMessage, &m); err != nil {
			return err
		}
		ed.Message = m
	default:
		return errors.New("unsupported event type during unmarshaling")
	}
	return nil
}

// Validate ensures the event has a recognized type and source
func (e *WsEvent) Validate() error {
	if e.Details.Type < MESSAGE || e.Details.Type > STATUS_UPDATE {
		return errors.New("invalid event type")
	}
	if e.Source < SERVER || e.Source > CLIENT {
		return errors.New("invalid event source")
	}
	if e.Details.Message == nil {
		return errors.New("missing event message payload")
	}
	return nil
}

// Marshal converts the WsEvent into a JSON representation
func (e *WsEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}
