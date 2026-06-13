package events

import (
	"encoding/json"
	"time"
)

type StatusMessage struct {
	UserID   int64     `json:"user_id"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

func (m StatusMessage) GetFromID() int64 { return m.UserID }
func (m StatusMessage) GetToID() int64   { return 0 }

func (m StatusMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *StatusMessage) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}
