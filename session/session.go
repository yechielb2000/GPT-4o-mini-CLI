package session

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"
)

const ChatCommands string = `Send and receive message in the terminal.
To close the conversation you can Ctrl+C or type #>exit.
To print history of messages type #>history.\n\n`

type Message struct {
	sender  string
	content []byte
}

type Session struct {
	sendChannel chan []byte
	conn        *websocket.Conn
	wg          sync.WaitGroup
	done        chan struct{}
}

func NewSession(u url.URL, apiKey string) (*Session, error) {
	headers := http.Header{}
	//TODO: I need to hold the apiKey in a better place right now ill just read from env
	headers.Add("Authorization", "Bearer "+apiKey)
	headers.Add("OpenAI-Beta", "realtime=v1")
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return nil, err
	}
	return &Session{
		sendChannel: make(chan []byte),
		conn:        conn,
		done:        make(chan struct{}),
	}, nil
}

func (s *Session) Start() {
	s.wg.Add(2)
	defer s.wg.Wait()

	go s.readData()
	go s.writeData()

	go func() {
		reader := bufio.NewScanner(os.Stdin)
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		for {
			fmt.Print("Type: ")
			if !reader.Scan() {
				break
			}
			text := reader.Text()

			if text == "#>exit" {
				log.Println("closing session...")
				s.Close()
				return
			}

			s.sendRequest(text)

			select {
			case <-s.done:
				log.Println("connection closed by server")
				return
			case <-interrupt:
				log.Println("interrupt received, closing connection")
				s.Close()
				return
			default:
			}
		}
	}()
}

func (s *Session) Close() {
	_ = s.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	time.Sleep(100 * time.Millisecond)
	s.conn.Close()
	close(s.sendChannel)
}

func (s *Session) readData() {
	defer close(s.done)

	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}
		s.handleMessage(message)
	}
}

func (s *Session) writeData() {
	for msg := range s.sendChannel {
		if err := s.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("write error:", err)
			return
		}
	}
}

func (s *Session) handleMessage(message []byte) {
	// TODO: extract and print the message
	fmt.Println("message:", string(message))
}

func (s *Session) sendRequest(text string) {
	// TODO: build the request and send it
	fmt.Println("sending:", text)
}
