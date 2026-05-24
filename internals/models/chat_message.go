package models

// {
//   sessionId: random string,
//   timestamp: epoch,
//   from: {id: UserId},
//   to: {id: UserId}
//   message: {"text":Text, attachments: [{type: "url", value: "https://"}, {"type":"imahge",value: "https:///"}],
//   token: ChatToken
// }

type Attachment struct{
	Type string `json:"type"`
	Value string `json:"value"`
}

type ChatToken string

type MessageContent struct{
	Text string `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type UserRef struct{
	Id int64 `json:"id"`
}

type ChatMessage struct{
	SessionId string `json:"sessionId"`
	Timestamp int64 `json:"timestamp"`
	From UserRef `json:"from"`
	To UserRef `json:"to"`
	Message MessageContent `json:"message"`
	Token ChatToken `json:"token"`
}