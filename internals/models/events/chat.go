package events

type Attachment struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ChatToken string

type MessageContent struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type ChatMessage struct {
	SessionId string         `json:"sessionId"`
	Timestamp int64          `json:"timestamp"`
	From      UserRef        `json:"from"`
	To        UserRef        `json:"to"`
	Message   MessageContent `json:"message"`
	Token     ChatToken      `json:"token"`
}

func (m ChatMessage) GetFromID() int64 { return m.From.Id }
func (m ChatMessage) GetToID() int64   { return m.To.Id }
