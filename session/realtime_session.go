package session

import (
	"bufio"
	"bytes"
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
	sendChannel   chan []byte
	conn          *websocket.Conn
	done          chan struct{}
	readyForInput chan struct{}
}

func NewRealtimeSession() (*RealtimeSession, error) {
	session := &RealtimeSession{
		sendChannel:   make(chan []byte),
		done:          make(chan struct{}),
		readyForInput: make(chan struct{}, 1),
	}
	session.Type = "realtime"

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

	go s.receiveData()
	go s.sendData()
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

		s.sendMessage(text)
	}
}

func (s *RealtimeSession) receiveData() {
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

func (s *RealtimeSession) sendData() {
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
	case events.ResponseCreated:
		fmt.Println(sessionRes)
	case events.ResponseTextDelta:
		fmt.Print(sessionRes["delta"])
	case events.ResponseDone:
		fmt.Println()
		s.readyForInput <- struct{}{}
	case events.Error:
		fmt.Println(sessionRes["error"])
	default:
		fmt.Println(sessionRes["type"])
	}
}

func (s *RealtimeSession) sendMessage(text string) {
	rawMessage, err := json.Marshal(builders.NewClientTextMessage(text))

	g := make(map[string]interface{})
	_ = json.Unmarshal(rawMessage, &g)
	fmt.Println(g)
	if err != nil {
		log.Println("marshal error:", err)
	}
	s.sendChannel <- rawMessage
}

func (s *RealtimeSession) handleFunctionCalls() {
	// TODO: maybe we will have an Action / Tool singleton with these options
}
