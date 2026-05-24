package ws

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

type Identifier string

type WsConnection struct{
	ID Identifier
	Conn *websocket.Conn
}

var wsUpgrader *websocket.Upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func NewConnection(id Identifier, w http.ResponseWriter, r *http.Request) (*WsConnection, error) {
	wsconn := WsConnection{ID: id} 
	conn, err := wsUpgrader.Upgrade(w,r,nil)
	if err != nil{
		return nil,errors.New("upgrade request failed")
	}
	wsconn.Conn = conn
	
	return &wsconn, nil
}

