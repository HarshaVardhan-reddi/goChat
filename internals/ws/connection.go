package ws

import (
	"chatonetoone/internals/models"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Identifier string

type WsConnection struct{
	ID Identifier
	Conn *websocket.Conn
	Broker chan models.WsEvent
}

var wsUpgrader *websocket.Upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: originCheck}

func NewConnection(id Identifier, w http.ResponseWriter, r *http.Request) (*WsConnection, error) {
	wsconn := WsConnection{ID: id, Broker: make(chan models.WsEvent, 1024)}
	conn, err := wsUpgrader.Upgrade(w,r,nil)
	if err != nil{
		return nil,errors.New("upgrade request failed")
	}
	wsconn.Conn = conn
	
	return &wsconn, nil
}

func originCheck(r *http.Request) bool{
	if os.Getenv("ENVIRONMENT") == "development"{
		return true
	}
	return os.Getenv("ALLOWED_ORIGIN") == r.Header.Get("Origin")
}