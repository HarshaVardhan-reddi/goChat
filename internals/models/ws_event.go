package models

import (
	"chatonetoone/internals/models/events"
	"encoding/json"
	"errors"
	"time"
)

type EventDetails struct {
	Type    events.EventType `json:"type"`
	Message events.Message   `json:"message"`
}

type WsEvent struct {
	From      events.UserRef     `json:"from"`
	To        events.UserRef     `json:"to"`
	Source    events.EventSource `json:"src"`
	Details   EventDetails       `json:"details"`
	Timestamp time.Time          `json:"timestamp"`
}

func NewWsEvent(fromId int, toId int, s events.EventSource, d json.RawMessage) (*WsEvent, error) {
	ed := EventDetails{}
	if err := json.Unmarshal(d,&ed); err != nil{
		return nil, err
	}

	return &WsEvent{
		From: events.UserRef{Id: int64(fromId)},
		To: events.UserRef{Id: int64(toId)},
		Source: s,
		Details: ed,
		Timestamp: time.Now(),
	}, nil
}

func (ed *EventDetails) UnmarshalJSON(data []byte) error {
	var temp struct {
		Type    events.EventType `json:"type"`
		Payload json.RawMessage  `json:"message"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	ed.Type = temp.Type

	switch ed.Type {
	case events.MESSAGE:
		ed.Message = &events.ChatMessage{}
	case events.SUBSCRIBE:
		ed.Message = &events.SubscribeMessage{}
	case events.UNSUBSCRIBE:
		ed.Message = &events.UnsubscribeMessage{}
	case events.STATUS_UPDATE:
		ed.Message = &events.StatusMessage{}
	default:
		return errors.New("unsupported event type")
	}

	return ed.Message.FromJSON(temp.Payload)
}

// MarshalJSON handles the polymorphic Message by calling ToJSON explicitly
func (ed EventDetails) MarshalJSON() ([]byte, error) {
	payload, err := ed.Message.ToJSON()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Type    events.EventType `json:"type"`
		Message json.RawMessage  `json:"message"`
	}{
		Type:    ed.Type,
		Message: payload,
	})
}

func (e *WsEvent) Validate() error {
	if e.Details.Type < events.MESSAGE || e.Details.Type > events.STATUS_UPDATE {
		return errors.New("invalid event type")
	}
	if e.Source < events.SERVER || e.Source > events.CLIENT {
		return errors.New("invalid event source")
	}
	return nil
}

func (e WsEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
