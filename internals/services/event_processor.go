package services

import (
	"chatonetoone/internals/models"
	"chatonetoone/internals/models/events"
	"chatonetoone/internals/ws"
	"fmt"
	"log"
	"strconv"
	"time"
)

// EventHandler defines the signature for functions that process specific event types
type EventHandler func(event models.WsEvent)

type EventProcessor struct {
	con      *ws.WsConnection
	hub      *ws.Hub
	handlers map[models.EventType]EventHandler
}

func NewEventProcessor(con *ws.WsConnection, hub *ws.Hub) *EventProcessor {
	ep := &EventProcessor{
		con:      con,
		hub:      hub,
		handlers: make(map[models.EventType]EventHandler),
	}

	ep.handlers[models.MESSAGE] = ep.handleChatMessage
	ep.handlers[models.SUBSCRIBE] = ep.handleSubscribe
	ep.handlers[models.UNSUBSCRIBE] = ep.handleUnsubscribe

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
	listener, err := redis.Subscribe(targetId)
	if err != nil {
		log.Println("Redis subscription error:", err)
		return
	}

	go func() {
		for {
			status, ok := <-listener
			if !ok || status == -1 {
				return
			}

			statusStr := "offline"
			if status == 1 {
				statusStr = "online"
			}

			// Using the NEW polymorphic structure from the events package
			statusPayload := events.StatusMessage{
				UserID:   event.Details.Message.GetToID(),
				Status:   statusStr,
				LastSeen: time.Now(),
			}

			resp := models.WsEvent{
				From:   events.UserRef{Id: event.Details.Message.GetToID()},
				To:     events.UserRef{Id: event.From.Id},
				Source: models.SERVER,
				Details: models.EventDetails{
					Type:    models.STATUS_UPDATE,
					Message: statusPayload,
				},
				Timestamp: time.Now(),
			}
			
			ep.con.Broker <- resp
		}
	}()
}

func (ep *EventProcessor) handleUnsubscribe(event models.WsEvent) {
	targetId := event.Details.Message.GetToID()
	fmt.Printf("User %d unsubscribed from %d\n", event.From.Id, targetId)
}
