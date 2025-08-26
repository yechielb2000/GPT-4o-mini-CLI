package session

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gpt4omini/api"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	config = api.GetConfig()
)

type Message struct {
	sender  string
	content []byte
}

type RealtimeSession struct {
	sendChannel    chan []byte
	conn           *websocket.Conn
	done           chan struct{}
	readyForInput  chan struct{}
	ConversationID string
	SessionID      string
	ClientSecret   ClientSecret
}

func NewRealtimeSession() (*RealtimeSession, error) {
	session := &RealtimeSession{
		sendChannel:   make(chan []byte),
		done:          make(chan struct{}),
		readyForInput: make(chan struct{}, 1),
	}
	createSessionRes, err := createNewSession()
	if err != nil {
		return nil, err
	}
	session.ClientSecret = createSessionRes.ClientSecret
	session.SessionID = createSessionRes.Id
	conn, err := session.createNewSocketConnection()
	if err != nil {
		return nil, err
	}
	session.conn = conn
	return session, err
}

func createNewSession() (*CreateSessionHTTPResponse, error) {
	bodyBytes, _ := json.Marshal(CreateSessionHTTPRequest{
		Modalities:   []string{"text"},
		Model:        config.Model.Name,
		Instructions: config.Model.Instruction,
	})
	u := "https://" + config.Api.Host + api.RealtimeSessionsPath
	req, _ := http.NewRequest("POST", u, bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+config.Api.Key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	sessionMetadata := &CreateSessionHTTPResponse{}
	if res.StatusCode == 200 {
		data, _ := io.ReadAll(res.Body)
		err = json.Unmarshal(data, &sessionMetadata)
	} else {
		err = errors.New("unexpected status code " + strconv.Itoa(res.StatusCode))
	}
	return sessionMetadata, err
}

func (s *RealtimeSession) createNewSocketConnection() (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+s.ClientSecret.Value)
	headers.Add("OpenAI-Beta", "realtime=v1")

	url := api.GetURL(api.RealtimePath)
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), headers)
	return conn, err
}

func (s *RealtimeSession) Start() {
	if s.HasExpired() {
		log.Println("RealtimeSession is expired, start a new one.")
		return
	}
	go s.readData()
	go s.writeData()
	go s.handleUserInput()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	s.readyForInput <- struct{}{}
	select {
	case <-s.done:
		log.Println("connection closed by server")
	case <-interrupt:
		log.Println("interrupt received, closing connection")
		s.Close()
	}
}

func (s *RealtimeSession) Close() {
	_ = s.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	time.Sleep(100 * time.Millisecond)
	s.conn.Close()
	close(s.sendChannel)
}

func (s *RealtimeSession) HasExpired() bool {
	expireTime := time.Unix(s.ClientSecret.ExpiresAt, 0)
	return time.Now().After(expireTime)
}

func (s *RealtimeSession) handleUserInput() {
	reader := bufio.NewScanner(os.Stdin)
	for range s.readyForInput {
		fmt.Print("#> ")
		if !reader.Scan() {
			log.Println("stdin closed, exiting input loop")
			s.Close()
			return
		}
		text := reader.Text()

		if text == "#>exit" {
			log.Println("closing session...")
			s.Close()
			return
		}

		s.sendRequest(text)
	}
}

func (s *RealtimeSession) readData() {
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

func (s *RealtimeSession) writeData() {
	for msg := range s.sendChannel {
		if err := s.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("write error:", err)
			return
		}
	}
}

func (s *RealtimeSession) handleMessage(message []byte) {
	sessionRes := map[string]any{}
	if err := json.Unmarshal(message, &sessionRes); err != nil {
		log.Println("unmarshal error:", err)
	}
	switch sessionRes["type"] {
	case "response.text.delta":
		fmt.Print(sessionRes["delta"])
	case "response.done":
		fmt.Println()
		s.readyForInput <- struct{}{}
	}
}

func (s *RealtimeSession) sendRequest(text string) {
	msgResConfig := MessageResponseConfig{
		Modalities:   []string{"text"},
		Instructions: text,
	}
	if s.ConversationID != "" {
		msgResConfig.Conversation = s.ConversationID
	}

	req := &MessageRequest{
		Type:     "response.create",
		Response: msgResConfig,
	}

	data, _ := json.Marshal(req)
	s.sendChannel <- data
}

func (s *RealtimeSession) handleFunctionCalls() {
	// TODO: maybe we will have an Action / Tool singleton with these options
}
