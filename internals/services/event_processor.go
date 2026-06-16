package services

import (
	"chatonetoone/internals/models"
	"chatonetoone/internals/models/events"
	"chatonetoone/internals/ws"
	"log"
	"strconv"
	"time"
)

// EventHandler defines the signature for functions that process specific event types
type EventHandler func(event models.WsEvent)

type EventProcessor struct {
	con      *ws.WsConnection
	hub      *ws.Hub
	handlers map[events.EventType]EventHandler
}

func NewEventProcessor(con *ws.WsConnection, hub *ws.Hub) *EventProcessor {
	ep := &EventProcessor{
		con:      con,
		hub:      hub,
		handlers: make(map[events.EventType]EventHandler),
	}

	ep.handlers[events.MESSAGE] = ep.handleChatMessage
	ep.handlers[events.SUBSCRIBE] = ep.handleSubscribe
	ep.handlers[events.UNSUBSCRIBE] = ep.handleUnsubscribe

	return ep
}

func (ep *EventProcessor) Process(event models.WsEvent) {
	if err := event.Validate(); err != nil {
		log.Println("Event validation failed:", err)
		return
	}

	if handler, ok := ep.handlers[event.Details.Type]; ok {
		handler(event)
	} else {
		log.Printf("No handler registered for type %d\n", event.Details.Type)
	}
}

func (ep *EventProcessor) handleChatMessage(event models.WsEvent) {
	targetId := strconv.FormatInt(event.To.Id, 10)
	targetCon, err := ep.hub.FetchConnection(ws.Identifier(targetId))
	if err == nil {
		targetCon.Broker <- event
	} else {
		log.Printf("User %s is offline.\n", targetId)
	}
}

func (ep *EventProcessor) handleSubscribe(event models.WsEvent) {
	targetId := strconv.FormatInt(event.Details.Message.GetToID(), 10)

	redis := FetchRedisConnection()
	listener, err := redis.Subscribe(ep.con.Ctx, targetId)
	if err != nil {
		log.Println("Redis subscription error:", err)
		return
	}

	go func() {
		for {
			select {
			case <-ep.con.Ctx.Done():
				return
			case status, ok := <-listener:
				if !ok {
					return
				}
				statusPayload := events.StatusMessage{
					UserID:   event.Details.Message.GetToID(),
					Status:   events.Status(status),
					LastSeen: time.Now(),
				}
				resp := models.WsEvent{
					From:   events.UserRef{Id: event.Details.Message.GetToID()},
					To:     events.UserRef{Id: event.From.Id},
					Source: events.SERVER,
					Details: models.EventDetails{
						Type:    events.STATUS_UPDATE,
						Message: &statusPayload,
					},
					Timestamp: time.Now(),
				}
				select {
				case ep.con.Broker <- resp:
				case <-ep.con.Ctx.Done():
					return
				}
			}
		}
	}()
}

func (ep *EventProcessor) handleUnsubscribe(event models.WsEvent) {
	log.Printf("User %d unsubscribed from %d\n", event.From.Id, event.Details.Message.GetToID())
}
