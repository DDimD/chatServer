package chat

// Message struct for sending message to clients
type OutputMessage struct {
	UserName string `json:"userName"`
	Body     string `json:"messageBody"`
}

type InputMessage struct {
	Body string `json:"messageBody"`
}
