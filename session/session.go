package session

type Message struct {
	sender  string
	content string
}

type Session struct {
	rcvMsgContent   chan []byte
	messagesHistory []Message
}

func NewSessionConfig() *Session {
	return &Session{
		rcvMsgContent:   make(chan []byte),
		messagesHistory: make([]Message, 0),
	}
}

func (cfg *Session) SendMessage() {
	// on user demand.
	// send a request containing the message.
	// the request returns and fills the rcvMsgContent channel...
}

func (cfg *Session) ReceiveMessage() {
	// print message in a streaming format
	// this function is a blocker
	// it gets message and as log as the channel full it will print
}

func (cfg *Session) Close() {
	// close session (the socket? we will see)
}

func (cfg *Session) ReturnToCLI() {
	// return to cli without losing this session
}

func (cfg *Session) saveMessage(sender string, content string) {
	cfg.messagesHistory = append(cfg.messagesHistory, Message{sender: sender, content: content})
}
