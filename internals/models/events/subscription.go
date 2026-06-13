package events

import "encoding/json"

type SubscribeMessage struct {
	TargetID int64 `json:"target_id"`
}

func (m SubscribeMessage) GetFromID() int64 { return 0 }
func (m SubscribeMessage) GetToID() int64   { return m.TargetID }

func (m SubscribeMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *SubscribeMessage) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}

type UnsubscribeMessage struct {
	TargetID int64 `json:"target_id"`
}

func (m UnsubscribeMessage) GetFromID() int64 { return 0 }
func (m UnsubscribeMessage) GetToID() int64   { return m.TargetID }

func (m UnsubscribeMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *UnsubscribeMessage) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}
