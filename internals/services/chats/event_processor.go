package chats

import (
	"chatonetoone/internals/models"
	"chatonetoone/internals/ws"
	"fmt"
	"log"
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

// Process is the main entry point
func (ep *EventProcessor) Process(event models.WsEvent) {

	if err := event.Validate(); err != nil {
		log.Println("Event validation failed for user", ep.con.ID, ":", err)
		return
	}

	if handler, ok := ep.handlers[event.EventType]; ok {
		handler(event)
	} else {
		log.Printf("No handler registered for event type %d from user %s\n", event.EventType, ep.con.ID)
	}
}

func (ep *EventProcessor) handleChatMessage(event models.WsEvent) {
	targetId, err := event.FetchTargetId()
	if err != nil {
		log.Println("Error fetching target ID for message from", ep.con.ID, ":", err)
		return
	}

	targetCon, err := ep.hub.FetchConnection(ws.Identifier(targetId))
	if err == nil {
		targetCon.Broker <- event
	} else {
		log.Printf("User %s is offline. Event should be saved to DB.\n", targetId)
	}
}

func (ep *EventProcessor) handleSubscribe(event models.WsEvent) {
	targetId, err := event.FetchTargetId()
	if err != nil {
		log.Println("Error fetching target ID for subscription from", ep.con.ID, ":", err)
		return
	}

	// Future: Implement Redis Pub/Sub status tracking
	fmt.Printf("User %s subscribed to status of %s\n", ep.con.ID, targetId)
}

func (ep *EventProcessor) handleUnsubscribe(event models.WsEvent) {
	targetId, err := event.FetchTargetId()
	if err != nil {
		log.Println("Error fetching target ID for unsubscription from", ep.con.ID, ":", err)
		return
	}

	fmt.Printf("User %s unsubscribed from status of %s\n", ep.con.ID, targetId)
}
