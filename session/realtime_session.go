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
	readyForInput    chan struct{}
	conn             *websocket.Conn
}

func NewRealtimeSession() (*RealtimeSession, error) {
	session := &RealtimeSession{
		BaseSession: BaseSession{
			Type: "realtime",
		},
		outgoingMessages: make(chan []byte),
		incomingMessages: make(chan []byte),
		messageChannel:   make(chan []byte),
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

	s.ctx, s.cancel = context.WithCancel(context.Background())
	defer s.cancel()

	go s.readMessages()
	go s.sendMessages()
	go s.handleIncomingMessage()
	go s.handleUserInput()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	s.readyForInput <- struct{}{}

	select {
	case <-interrupt:
		log.Println("Interrupt received, closing connection")
		s.Close()
	}
}

func (s *RealtimeSession) Close() {
	if s.conn != nil {
		_ = s.conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		time.Sleep(100 * time.Millisecond)
		if err := s.conn.Close(); err != nil {
			log.Println("Warning: Failed to close connection", err)
		}
	}

	if s.cancel != nil {
		s.cancel()
	}

	close(s.outgoingMessages)
	close(s.incomingMessages)
	close(s.readyForInput)

}

func (s *RealtimeSession) handleUserInput() {
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
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *RealtimeSession) readMessages() {
	for {
		select {
		case <-s.ctx.Done():
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

func (s *RealtimeSession) sendMessages() {
	for {
		select {
		case message := <-s.outgoingMessages:
			if err := s.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("send error:", err)
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *RealtimeSession) handleIncomingMessage() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case message := <-s.incomingMessages:
			sessionRes := map[string]any{}
			if err := json.Unmarshal(message, &sessionRes); err != nil {
				log.Println("unmarshal error:", err)
				return
			}
			switch sessionRes["type"] {
			case events.ResponseDone:
				fmt.Println()
				select {
				case s.readyForInput <- struct{}{}:
				default:
				}
			case events.ResponseTextDelta:
				fmt.Print(sessionRes["delta"])
			case events.Error:
				if errObj, ok := sessionRes["error"].(map[string]any); ok {
					fmt.Println(errObj["message"])
				}
			}
		}
	}
}

func (s *RealtimeSession) handleFunctionCalls() {
	// TODO: maybe we will have an Action / Tool singleton with these options
}
