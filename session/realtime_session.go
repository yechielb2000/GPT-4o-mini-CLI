package session

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gpt4omini/config"
	"gpt4omini/global_tools"
	"gpt4omini/types"
	"log"
	"net/http"
	"os"
	"time"
)

type RealtimeSession struct {
	BaseSession

	conn *websocket.Conn
}

func NewRealtimeSession() (*RealtimeSession, error) {
	session := &RealtimeSession{
		BaseSession: BaseSession{
			Type:             "realtime",
			outgoingMessages: make(chan types.ConversationItem),
			functionCalls:    make(chan types.ConversationItem),
			incomingEvents:   make(chan types.Event),
			readyForInput:    make(chan struct{}, 1),
		},
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
	s.ctx, s.cancel = context.WithCancel(context.Background())
	defer s.cancel()

	s.wg.Add(5)
	go s.readMessages()
	go s.sendMessages()
	go s.handleIncomingEvents()
	go s.handleFunctionCalls()
	go s.handleUserInput()

	s.readyForInput <- struct{}{}

	select {
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
	close(s.incomingEvents)
	close(s.functionCalls)
	close(s.readyForInput)
}

func (s *RealtimeSession) establishConnection() (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+s.clientSecret.Value)
	headers.Add("OpenAI-Beta", "realtime=v1")

	url := config.GetURL(config.RealtimePath)
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), headers)
	return conn, err
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
			s.outgoingMessages <- types.NewClientTextConversationItem(text)
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
			event := types.Event{}
			if err := json.Unmarshal(msg, &event); err != nil {
				log.Println("unmarshal error:", err)
				return
			}
			if (event.Item.Type == types.FunctionCallItem || event.Item.Type == types.MessageItem) &&
				event.Item.Status == types.Completed {
				s.AddToConversation(event.Item)
			}
			s.incomingEvents <- event
		}
	}
}

func (s *RealtimeSession) sendMessages() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.outgoingMessages:
			s.AddToConversation(msg)
			message, err := json.Marshal(s.NewClientMessage())
			if err != nil {
				log.Println("marshal error:", err)
				continue
			}
			if err = s.conn.WriteMessage(websocket.TextMessage, message); err != nil {
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
		case event := <-s.incomingEvents:
			switch event.Type {
			case types.ResponseDoneEvent:
				fmt.Println()
				s.readyForInput <- struct{}{}
			case types.ResponseTextDeltaEvent:
				for _, r := range event.Delta {
					fmt.Printf("%c", r)
					time.Sleep(22 * time.Millisecond)
				}
			case types.ResponseOutputItemDoneEvent:
				if event.Item.Type == types.FunctionCallItem {
					s.functionCalls <- event.Item
				}
			case types.ErrorEvent:
				fmt.Println("error message:", event.Error.Message)
			}
		}
	}
}

func (s *RealtimeSession) handleFunctionCalls() {

	for {
		select {
		case <-s.ctx.Done():
			return
		case item := <-s.functionCalls:
			if item.Name == ExitSessionFunctionName {
				log.Println("Closing session...")
				s.close()
			}

			arguments, err := item.GetArguments()
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			result, err := global_tools.CallFunction(item.Name, arguments)
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			toolResItem := types.NewClientFunctionCallConversationItem(item, result)
			s.outgoingMessages <- toolResItem
		}
	}
}
