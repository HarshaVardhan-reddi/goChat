package events

type SubscribeMessage struct {
	TargetID int64 `json:"target_id"`
}

func (m SubscribeMessage) GetFromID() int64 { return 0 }
func (m SubscribeMessage) GetToID() int64   { return m.TargetID }

type UnsubscribeMessage struct {
	TargetID int64 `json:"target_id"`
}

func (m UnsubscribeMessage) GetFromID() int64 { return 0 }
func (m UnsubscribeMessage) GetToID() int64   { return m.TargetID }
