package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
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
