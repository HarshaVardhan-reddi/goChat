package events

import (
	"encoding/json"
	"time"
)

type Status int
const(
	INACTIVE = iota + 0
	ACTIVE
)

type StatusMessage struct {
	UserID   int64     `json:"user_id"`
	Status   Status    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

func NewEvent(uid int64, ls time.Time, s Status) StatusMessage {
	return  StatusMessage{UserID: uid, LastSeen: ls, Status: s}
}

func (m StatusMessage) GetFromID() int64 { return m.UserID }
func (m StatusMessage) GetToID() int64   { return 0 }

func (m StatusMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *StatusMessage) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}
