package events

// UserRef represents a minimal user reference
type UserRef struct {
	Id int64 `json:"id"`
}

type EventType int

const (
	MESSAGE EventType = iota + 1
	SUBSCRIBE
	UNSUBSCRIBE
	STATUS_UPDATE
)

type EventSource int

const (
	SERVER EventSource = iota + 1
	CLIENT
)

// Message is the interface for all event payloads.
// It uses explicit ToJSON/FromJSON methods to avoid recursion issues.
type Message interface {
	GetFromID() int64
	GetToID() int64
	ToJSON() ([]byte, error)
	FromJSON(data []byte) error
}
