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

type Session struct {
	sendChannel    chan []byte
	conn           *websocket.Conn
	done           chan struct{}
	ConversationID string
	SessionID      string
	ClientSecret   ClientSecret
}

func NewSession() (*Session, error) {
	session := &Session{
		sendChannel: make(chan []byte),
		done:        make(chan struct{}),
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

func (s *Session) createNewSocketConnection() (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+s.ClientSecret.Value)
	headers.Add("OpenAI-Beta", "realtime=v1")

	url := api.GetURL(api.RealtimePath)
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), headers)
	return conn, err
}

func (s *Session) Start() {
	if s.HasExpired() {
		log.Println("Session is expired, start a new one.")
		return
	}
	go s.readData()
	go s.writeData()

	reader := bufio.NewScanner(os.Stdin)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		//TODO: delete it and make this actions using function calls.
		fmt.Print("#>")
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

func (s *Session) HasExpired() bool {
	expireTime := time.Unix(s.ClientSecret.ExpiresAt, 0)
	return time.Now().After(expireTime)
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
	fmt.Println("res:", string(message))
}

func (s *Session) sendRequest(text string) {
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

func (s *Session) handleFunctionCalls() {
	// TODO: maybe we will have an Action / Tool singleton with these options
}
