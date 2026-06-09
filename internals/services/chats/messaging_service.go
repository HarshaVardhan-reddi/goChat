package chats

import (
	"chatonetoone/internals/models"
	"chatonetoone/internals/ws"

	"encoding/json"
	"log"
)

type MessagingService struct {
	Connection *ws.WsConnection
	hub        *ws.Hub
}

func StartNewMessagingService(con *ws.WsConnection, hub *ws.Hub) {
	msgservice := MessagingService{Connection: con, hub: hub}

	go msgservice.MessageReader()
	go msgservice.MessageWriter()
}

func (ms *MessagingService) MessageReader() {
	processor := NewEventProcessor(ms.Connection, ms.hub)

	defer func() {
		ms.Connection.Conn.Close()
		ms.hub.DeleteConnection(ms.Connection.ID)
		close(ms.Connection.Broker)
	}()

	for {
		_, rawMessage, err := ms.Connection.Conn.ReadMessage()
		if err != nil {
			log.Println("error while reading the data from connection for id:", ms.Connection.ID, err)
			return
		}

		wse := models.WsEvent{}
		if err := json.Unmarshal(rawMessage, &wse); err != nil {
			log.Println("error unmarshaling event from id:", ms.Connection.ID, err)
			continue
		}

		processor.Process(wse)
	}
}

func (ms *MessagingService) MessageWriter() {
	defer func() {
		ms.hub.DeleteConnection(ms.Connection.ID)
		ms.Connection.Conn.Close()
	}()

	for {
		msg, ok := <-ms.Connection.Broker
		if !ok {
			log.Println("channel closed for id:", ms.Connection.ID)
			return
		}

		message, err := json.Marshal(msg)
		if err != nil {
			log.Println("error encoding message for id:", ms.Connection.ID, err)
			continue
		}

		if err := ms.Connection.Conn.WriteMessage(1, message); err != nil {
			log.Println("error sending message to id:", ms.Connection.ID, err)
			return
		}
	}
}
