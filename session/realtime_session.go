package session

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gpt4omini/builders"
	"gpt4omini/config"
	"gpt4omini/events"
	"gpt4omini/global_tools"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	if createSessionRes, err := ConfigureModel(); err == nil {
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

func (s *RealtimeSession) Start() {
	if s.HasClientSecretExpired() {
		log.Println("The session is expired, start a new one.")
		return
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	defer s.cancel()

	s.wg.Add(4)
	go s.readMessages()
	go s.sendMessages()
	go s.handleIncomingEvents()
	go s.handleUserInput()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	s.readyForInput <- struct{}{}

	select {
	case <-interrupt:
		log.Println("Interrupt received, closing connection")
		s.close()
	case <-s.ctx.Done():
	}
}

func (s *RealtimeSession) close() {
	if s.cancel != nil {
		s.cancel()
	}

	if s.conn != nil {
		_ = s.conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		_ = s.conn.Close()
	}

	s.wg.Wait()
	close(s.outgoingMessages)
	close(s.incomingMessages)
	close(s.readyForInput)
}

func (s *RealtimeSession) handleUserInput() {
	defer s.wg.Done()
	reader := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.readyForInput:
			fmt.Printf("(%s)> ", s.GetID())
			if !reader.Scan() {
				log.Println("reader closed")
				s.close()
				return
			}
			text := reader.Text()
			if text == "#>exit" {
				//TODO: make function call
				log.Println("closing session...")
				s.close()
				return
			}
			message, err := json.Marshal(builders.NewClientTextMessage(text))
			if err != nil {
				log.Println("marshal error:", err)
				continue
			}
			s.outgoingMessages <- message
		}
	}
}

func (s *RealtimeSession) readMessages() {
	defer s.wg.Done()
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
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case message := <-s.outgoingMessages:
			if err := s.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("send error:", err)
			}
		}
	}
}

func (s *RealtimeSession) handleIncomingEvents() {
	defer s.wg.Done()
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
				delta := sessionRes["delta"].(string)
				for _, r := range delta {
					fmt.Printf("%c", r)
					time.Sleep(22 * time.Millisecond)
				}
			case events.ResponseOutputItemDone:
				item := sessionRes["item"].(map[string]any)
				if item["type"] == "function_call" {
					name := item["name"].(string)
					var args map[string]interface{}
					switch v := item["arguments"].(type) {
					case string:
						if err := json.Unmarshal([]byte(v), &args); err != nil {
							log.Println("failed to unmarshal arguments:", err)
							return
						}
					case map[string]interface{}:
						args = v
					}
					function := global_tools.Factory[name]
					result := function(args)
					fmt.Println("result:", result)
				}
			case events.Error:
				if errObj, ok := sessionRes["error"].(map[string]any); ok {
					fmt.Println(errObj["message"])
				}
			}
		}
	}
}

func handleFunctionCalls() {

}

func (s *RealtimeSession) establishConnection() (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+s.clientSecret.Value)
	headers.Add("OpenAI-Beta", "realtime=v1")

	url := config.GetURL(config.RealtimePath)
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), headers)
	return conn, err
}
