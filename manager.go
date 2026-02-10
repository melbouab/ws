package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SenMessage
}

func SenMessage(event Event, c *Client) error {
	fmt.Println("")
	return nil
}

func (m *Manager) routEvent(event Event, c *Client) error {
	// check if the event type
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) servWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	Client := NewClient(conn, m)
	m.addClient(Client)

	// Start clinet process

	go Client.ReadMessage()
	go Client.writeMessages()
}

func (m *Manager) addClient(Client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[Client] = true
}

func (m *Manager) removeClient(Client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[Client]; ok {
		Client.connection.Close()
		delete(m.clients, Client)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}
