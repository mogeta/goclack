package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/net/websocket"
)

// Dice act roll dice
type Dice struct {
	channel string
}

// A Event is response type
type Event struct {
	Type string `json:"type"`
}

// A Message is response type by slack websocket
type Message struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// Send is Message Struct for slack websocket
type Send struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// NewDiceTask is
func NewDiceTask() Dice {
	d := Dice{}
	return d
}

// Check receive slack message
func (d *Dice) Check(b []byte) bool {
	var message Message
	err := json.Unmarshal(b, &message)
	if err != nil {
		return false
	}

	if message.Text == "dice" {
		d.channel = message.Channel
		return true
	}

	return false
}

// Action is
func (d *Dice) Action(ws *websocket.Conn) {
	r := rand.Intn(6) + 1

	var send Send
	send.ID = 1
	send.Type = "message"
	send.Channel = d.channel
	send.Text = fmt.Sprintf("dice roll : %d", r)

	websocket.JSON.Send(ws, send)
}

func init() {
	rand.Seed(time.Now().Unix())
}
