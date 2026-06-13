package events

// UserRef represents a minimal user reference
type UserRef struct {
	Id int64 `json:"id"`
}

// Message is the interface for all event payloads (Chat, Status, Subscription)
type Message interface {
	GetFromID() int64
	GetToID() int64
}
