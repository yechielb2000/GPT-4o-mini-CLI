package session

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gpt4omini/builders"
	"gpt4omini/config"
	"gpt4omini/events"
	"gpt4omini/types"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type RealtimeSession struct {
	BaseSession
	outgoingMessages chan []byte
	incomingMessages chan []byte
	messageChannel   chan []byte
	conn             *websocket.Conn
	done             chan struct{}
	readyForInput    chan struct{}
}

func NewRealtimeSession() (*RealtimeSession, error) {
	session := &RealtimeSession{
		BaseSession: BaseSession{
			Type: "realtime",
		},
		outgoingMessages: make(chan []byte),
		incomingMessages: make(chan []byte),
		messageChannel:   make(chan []byte),
		done:             make(chan struct{}),
		readyForInput:    make(chan struct{}, 1),
	}

	if createSessionRes, err := configureModel(); err == nil {
		session.ID = createSessionRes.Id
		session.clientSecret = createSessionRes.ClientSecret
		session.createdAt = time.Now()
	} else {
		return nil, err
	}

	if conn, err := session.establishConnection(); err == nil {
		session.conn = conn
	} else {
		return nil, err
	}
	return session, nil
}

func configureModel() (*types.ConfigureModelResponse, error) {
	bodyBytes, _ := json.Marshal(types.ConfigureModelRequest{
		Modalities:   []string{"text"},
		Model:        cfg.Model.Name,
		Instructions: cfg.Model.Instruction,
	})
	u := "https://" + cfg.Api.Host + config.RealtimeSessionsPath
	req, _ := http.NewRequest("POST", u, bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+cfg.Api.Key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	sessionMetadata := &types.ConfigureModelResponse{}
	if res.StatusCode == 200 {
		err = json.Unmarshal(body, &sessionMetadata)
	} else {
		err = errors.New("unexpected status code " + strconv.Itoa(res.StatusCode) + ".\n" + string(body))
	}
	return sessionMetadata, err
}

func (s *RealtimeSession) establishConnection() (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+s.clientSecret.Value)
	headers.Add("OpenAI-Beta", "realtime=v1")

	url := config.GetURL(config.RealtimePath)
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), headers)
	return conn, err
}

func (s *RealtimeSession) Start() {
	if s.HasClientSecretExpired() {
		log.Println("The session is expired, start a new one.")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.readMessages(ctx)
	go s.sendMessages(ctx)
	go s.handleUserInput(ctx)

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
}

func (s *RealtimeSession) handleUserInput(ctx context.Context) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-s.readyForInput:
			fmt.Printf("(%s)> ", s.GetID())
			if !reader.Scan() {
				log.Println("reader closed")
				s.Close()
				return
			}
			text := reader.Text()
			if text == "#>exit" {
				log.Println("closing session...")
				s.Close()
				return
			}
			message, err := json.Marshal(builders.NewClientTextMessage(text))
			if err != nil {
				log.Println("marshal error:", err)
				continue
			}
			s.outgoingMessages <- message
		case <-ctx.Done():
			return
		}
	}
}

func (s *RealtimeSession) readMessages(ctx context.Context) {
	defer close(s.done)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := s.conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			s.incomingMessages <- msg
		}
	}
}

func (s *RealtimeSession) sendMessages(ctx context.Context) {
	for {
		select {
		case message := <-s.outgoingMessages:
			if err := s.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("send error:", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *RealtimeSession) handleMessage(message []byte) {
	sessionRes := map[string]any{}
	if err := json.Unmarshal(message, &sessionRes); err != nil {
		log.Println("unmarshal error:", err)
		return
	}

	switch sessionRes["type"] {
	case events.ResponseCreated:
		log.Println("session created")
	case events.ResponseTextDelta:
		fmt.Println(sessionRes["delta"])
	case events.ResponseDone:
		select {
		case s.readyForInput <- struct{}{}:
		default:
		}
	case events.Error:
		if errObj, ok := sessionRes["error"].(map[string]any); ok {
			fmt.Println(errObj["message"])
		}
	}
}

func (s *RealtimeSession) handleFunctionCalls() {
	// TODO: maybe we will have an Action / Tool singleton with these options
}
