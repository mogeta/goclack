package goclack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

const (
	apiStart = "https://slack.com/api/rtm.start"
	origin   = "https://slack.com/"
)

// A Start is init data from https://api.slack.com/methods/rtm.start
type Start struct {
	Ok  bool   `json:"ok"`
	URL string `json:"url"`
}

// A Events is response type
type Events struct {
	Type string `json:"type"`
}

// Behavior is
// Condition...条件
// Action 条件を満たした場合の...処理
type Behavior interface {
	Check(b []byte) bool
	Action(ws *websocket.Conn)
}

// A Goclack is
type Goclack struct {
	token      string
	behaiviors []Behavior
}

// New is Create Goclack Struct
func New(token string) *Goclack {
	g := &Goclack{token: token}
	return g
}

// Run is ...
func (g *Goclack) Run() {
	param := g.start()
	ws := g.connect(param.URL)
	g.relay(ws)
}

// AddBehaiviors is
func (g *Goclack) AddBehaiviors(b Behavior) {
	g.behaiviors = append(g.behaiviors, b)
}

func ping(ws *websocket.Conn, ch chan<- error) {
	for id := 0; ; id++ {
		message := fmt.Sprintf(`{"id":%v, "type": "ping"}`, id)
		err := websocket.Message.Send(ws, message)
		if err != nil {
			ch <- err
			return
		}
		time.Sleep(10 * time.Second)
	}
}

func (g *Goclack) relay(ws *websocket.Conn) error {
	ch1 := make(chan error)
	ch2 := make(chan error)

	go g.receive(ws, ch1)
	go ping(ws, ch2)

	select {
	case err1 := <-ch1:
		return err1
	case err2 := <-ch2:
		return err2
	}
}

//https://api.slack.com/rtm
func (g *Goclack) connect(wsurl string) *websocket.Conn {
	//origin := "https://slack.com/"
	ws, err := websocket.Dial(wsurl, "", origin)
	if err != nil {
		panic(err)
	}
	return ws
}

// receive slack send message.
func (g *Goclack) receive(ws *websocket.Conn, ch chan<- error) {
	//var events Events
	var data []byte
	var event Events

	for {
		if err := websocket.Message.Receive(ws, &data); err != nil {
			ch <- err
			break
		}

		err := json.Unmarshal(data, &event)
		if err != nil {
			ch <- err
			break
		}
		fmt.Println(event.Type)

		for _, behaivior := range g.behaiviors {
			if behaivior.Check(data) {
				behaivior.Action(ws)
			}
		}

	}
}

//https://api.slack.com/methods/rtm.start
func (g *Goclack) start() Start {
	resp, err := http.Get(apiStart + "?&pretty=1&token=" + g.token)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var start Start
	err = json.Unmarshal(body, &start)
	if err != nil {
		panic(err)
	}
	return start

}
