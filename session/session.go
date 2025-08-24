package session

type Session struct {
	rcvMsg chan []byte
}

func NewSessionConfig() *Session {
	return &Session{
		rcvMsg: make(chan []byte),
	}
}

func (cfg *Session) SendMessage() {
	// on user demand.
	// send a request containing the message.
	// the request returns and fills the rcvMsg channel...
}

func (cfg *Session) ReceiveMessage() {
	// print message in a streaming format
	// this function is a blocker
	// it gets message and as log as the channel full it will print
}

func (cfg *Session) Close() {
	// close session (the socket? we will see)
}
