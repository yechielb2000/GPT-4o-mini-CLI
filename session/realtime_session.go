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
	"gpt4omini/types"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type RealtimeSession struct {
	BaseSession
	outgoingMessages chan types.ConversationItem
	incomingEvents   chan events.Event
	messageChannel   chan []byte
	readyForInput    chan struct{}
	conn             *websocket.Conn
}

func NewRealtimeSession() (*RealtimeSession, error) {
	session := &RealtimeSession{
		BaseSession: BaseSession{
			Type: "realtime",
		},
		outgoingMessages: make(chan types.ConversationItem),
		incomingEvents:   make(chan events.Event),
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
	close(s.incomingEvents)
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
			s.outgoingMessages <- builders.NewClientTextConversationItem(text)
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
			event := events.Event{}
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
			conversationItem := builders.NewClientConversationEvent(s.GetConversation())
			fmt.Println(conversationItem)
			message, err := json.Marshal(conversationItem)
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
			case events.ResponseDone:
				fmt.Println()
				s.readyForInput <- struct{}{}
			case events.ResponseTextDelta:
				for _, r := range event.Delta {
					fmt.Printf("%c", r)
					time.Sleep(22 * time.Millisecond)
				}
			case events.ResponseOutputItemDone:
				fmt.Println("item call id", event.Item.CallID)
				fmt.Println("item id", event.Item.ID)
				if event.Item.Type == types.FunctionCallItem {
					if err := s.handleFunctionCalls(event.Item); err != nil {
						log.Println("handleFunctionCalls error:", err)
						return
					}
				}
			case events.Error:
				fmt.Println("error message:", event.Error.Message)
			}
		}
	}
}

func (s *RealtimeSession) handleFunctionCalls(item types.ConversationItem) error {

	if item.Name == ExitSessionFunctionName {
		log.Println("Closing session...")
		s.close()
	}

	arguments, err := item.GetArguments()
	if err != nil {
		return err
	}

	result, err := global_tools.CallFunction(item.Name, arguments)
	if err != nil {
		return err
	}

	toolResItem := builders.NewClientFunctionCallConversationItem(item, result)
	s.outgoingMessages <- toolResItem
	return nil
}

func (s *RealtimeSession) establishConnection() (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+s.clientSecret.Value)
	headers.Add("OpenAI-Beta", "realtime=v1")

	url := config.GetURL(config.RealtimePath)
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), headers)
	return conn, err
}
