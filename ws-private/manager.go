package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

type Manager struct {
	clients  ClientList
	users    map[string]*Client // username -> Client
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		users:    make(map[string]*Client),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = m.SendMessage
	m.handlers[EventPrivateMessage] = m.SendPrivateMessage
	m.handlers[EventLogin] = m.LoginHandler
}

// Handler ديال Login
func (m *Manager) LoginHandler(event Event, c *Client) error {
	var loginEvent LoginEvent
	if err := json.Unmarshal(event.Payload, &loginEvent); err != nil {
		return err
	}

	m.Lock()
	defer m.Unlock()

	// تحقق إلا Username مستعمل
	if _, exists := m.users[loginEvent.Username]; exists {
		return errors.New("username already taken")
	}

	// سجل الـ username
	c.username = loginEvent.Username
	m.users[loginEvent.Username] = c

	log.Printf("User %s logged in", loginEvent.Username)
	return nil
}

// Handler ديال Public Message (broadcast لجميع الـ clients)
func (m *Manager) SendMessage(event Event, c *Client) error {
	var messageEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &messageEvent); err != nil {
		return err
	}

	// كرييي event جديد باش نسيفطوه
	newEvent := NewMessageEvent{
		Message: messageEvent.Message,
		From:    c.username,
		Sent:    time.Now().Format("15:04:05"),
	}

	data, err := json.Marshal(newEvent)
	if err != nil {
		return err
	}

	outgoingEvent := Event{
		Type:    EventNewMessage,
		Payload: data,
	}

	// سيفط لجميع الـ clients
	m.RLock()
	defer m.RUnlock()

	for client := range m.clients {
		client.egress <- outgoingEvent
	}

	log.Printf("Broadcast message from %s: %s", c.username, messageEvent.Message)
	return nil
}

// Handler ديال Private Message
func (m *Manager) SendPrivateMessage(event Event, c *Client) error {
	var pmEvent PrivateMessageEvent
	if err := json.Unmarshal(event.Payload, &pmEvent); err != nil {
		return err
	}

	// جيب الـ target client
	m.RLock()
	targetClient, ok := m.users[pmEvent.To]
	m.RUnlock()

	if !ok {
		log.Printf("User %s not found", pmEvent.To)
		return errors.New("user not found: " + pmEvent.To)
	}

	// كرييي event جديد
	newEvent := NewMessageEvent{
		Message: "[Private] " + pmEvent.Message,
		From:    c.username,
		Sent:    time.Now().Format("15:04:05"),
	}

	data, err := json.Marshal(newEvent)
	if err != nil {
		return err
	}

	outgoingEvent := Event{
		Type:    EventNewMessage,
		Payload: data,
	}

	// سيفط الرسالة للـ target فقط
	targetClient.egress <- outgoingEvent

	// سيفط نسخة للمرسل باش يشوف الرسالة ديالو
	c.egress <- outgoingEvent

	log.Printf("Private message from %s to %s: %s", c.username, pmEvent.To, pmEvent.Message)
	return nil
}

func (m *Manager) routEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	}
	return errors.New("there is no such event type: " + event.Type)
}

func (m *Manager) servWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	client := NewClient(conn, m)
	m.addClient(client)

	// Start client processes
	go client.ReadMessage()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
	log.Println("Client added, total clients:", len(m.clients))
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
		
		// حيد من users map
		if client.username != "" {
			delete(m.users, client.username)
			log.Printf("User %s disconnected", client.username)
		}
		
		log.Println("Client removed, total clients:", len(m.clients))
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:8080":
		return true
	default:
		return true // للتجربة فقط، في production خاصك تكون strict
	}
}
