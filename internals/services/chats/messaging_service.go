package chats

import (
	"chatonetoone/internals/models"
	"chatonetoone/internals/ws"

	"encoding/json"
	"strconv"

	"log"

	"github.com/gorilla/websocket"
)

type MessagingService struct {
	Connection *ws.WsConnection
	hub        *ws.Hub
}

func StartNewMessagingService(con *ws.WsConnection, hub *ws.Hub) {
	msgservice := MessagingService{Connection: con, hub: hub}

	go msgservice.MessageReaderWithContext()
	go msgservice.MessageWriterWithContext()
}

func (ms *MessagingService) MessageReaderWithContext() {
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
		chatmsg := models.ChatMessage{}
		json.Unmarshal(rawMessage, &chatmsg)
		id := strconv.Itoa(int(chatmsg.To.Id))
		TargetCon, errInFetch := ms.hub.FetchConnection(ws.Identifier(id))
		msgInTxt := chatmsg.Message.Text
		if errInFetch == nil {
			
			TargetCon.Broker <- msgInTxt
		}else{
			log.Println("sending push message here..")
			log.Println("this is you message", msgInTxt)
		}
	}
}

func (ms *MessagingService) MessageWriterWithContext() {
	defer func() {
		ms.hub.DeleteConnection(ms.Connection.ID)
		ms.Connection.Conn.Close()
	}()
	for {
		msg, ok := <-ms.Connection.Broker
		if !ok {
			log.Println("channel close in id:", ms.Connection.ID)
			return
		}
		if err := ms.Connection.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("err in sending msg: ", err, "id: ", ms.Connection.ID)
			return
		}
	}
}
