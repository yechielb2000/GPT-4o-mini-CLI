package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type Message struct {
	sender  string
	content []byte
}

type Session struct {
	sendChannel    chan []byte
	conn           *websocket.Conn
	done           chan struct{}
	conversationID string
}

func NewSession(u url.URL, apiKey string) (*Session, error) {
	fmt.Println(u.String())
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+apiKey)
	fmt.Println(apiKey)
	fmt.Printf("API key length: %d\n", len(apiKey))
	headers.Add("OpenAI-Beta", "realtime=v1")
	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			fmt.Println("HTTP Status:", resp.Status)
			fmt.Println("Body:", string(body))
		}
		log.Fatal(err)
	}
	return &Session{
		sendChannel: make(chan []byte),
		conn:        conn,
		done:        make(chan struct{}),
	}, nil
}

func (s *Session) Start() {
	go s.readData()
	go s.writeData()

	s.sendRequest(`You are a rock & roll biggest fan that can help me learn anything about The Scorpions`)

	reader := bufio.NewScanner(os.Stdin)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		for {
			//TODO: delete it and make this actions using function calls.
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
	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		log.Println("json parse error:", err)
		return
	}
	if data["type"] == "response.created" {
		if resp, ok := data["response"].(map[string]interface{}); ok {
			if id, ok := resp["conversation"].(string); ok {
				s.conversationID = id
			}
		}
	}

	switch data["type"] {
	case "response.output_text.delta":
		if delta, ok := data["delta"].(string); ok {
			fmt.Print(delta)
		}
	case "response.completed":
		fmt.Print("\n")
	}

}

func (s *Session) sendRequest(text string) {
	req := map[string]interface{}{
		"type": "response.create",
		"response": map[string]interface{}{
			"modalities":   []string{"text"},
			"instructions": text,
		},
	}
	if s.conversationID != "" {
		req["response"].(map[string]interface{})["conversation"] = s.conversationID
	}

	data, _ := json.Marshal(req)
	s.sendChannel <- data
}

func (s *Session) handleFunctionCalls() {
	// TODO: maybe we will have an Action / Tool singleton with these options
}

func LoadSessionConfigParams(data []byte) {
	// TODO: load from config.yaml the api info and session params
}
