package main

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage    = "send_message"
	EventPrivateMessage = "private_message"
	EventLogin          = "login"
	EventNewMessage     = "new_message"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type PrivateMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
	To      string `json:"to"`
}

type LoginEvent struct {
	Username string `json:"username"`
}

type NewMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
	Sent    string `json:"sent"`
}
