package ws

import (
	"errors"
	"sync"
)

var ChatHub *Hub = &Hub{}

type Hub struct{
	connections sync.Map // this ds wont suite for real world application as it is ready heavy not for heavy writes
}

func (h *Hub) PutConnection(cd *WsConnection) error {
	
	if cd.ID == ""{
		return errors.New("identifier is missing in the arguments")
	}

	if cd.Conn == nil{
		return errors.New("connection is missing in the arguments")
	}

	h.connections.Store(cd.ID, cd)
	return nil
}

func (h *Hub) FetchConnection(id Identifier) (*WsConnection,error) {
	val, ok := h.connections.Load(id)
	if !ok{
		return nil, errors.New("unable to find the connection")
	}
	conn, succ := val.(*WsConnection)
	if !succ{
		return nil, errors.New("stored value is not connection type")
	}
	return conn, nil
}