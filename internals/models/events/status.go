package events

import "time"

type StatusMessage struct {
	UserID   int64     `json:"user_id"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

func (m StatusMessage) GetFromID() int64 { return m.UserID }
func (m StatusMessage) GetToID() int64   { return 0 }
