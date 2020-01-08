package sockets

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"net"

	"github.com/ahmetcanozcan/reves/sockets/messages"
)

//Socket : an abstraction for handling tcp connection
type Socket struct {
	events      map[string]func(messages.Payload)
	id          string
	conn        *net.Conn
	initialized bool
}

//On : Dispatch an event handler to event listener
func (s *Socket) On(name string, handler func(messages.Payload)) {
	s.events[name] = handler
}

//Emit : Send payload to client
func (s *Socket) Emit(name string, payload messages.Payload) {
	message := fmt.Sprintf("%s;%s\n", name, payload.Compile())
	fmt.Fprintf(*s.conn, message)
}

//IsOk : returns the socket complete its pre-communicaion tasks
func (s *Socket) IsOk() bool {
	return s.initialized
}

//GetID :
func (s Socket) GetID() string {
	return s.id
}

//Join :
func (s *Socket) Join(name string) {
	GetRoom(name).AddSocket(s)
}

//JoinMatchMaking :
func (s *Socket) JoinMatchMaking() {
	GetMatchMakingRoom().AddSocket(s)
}

//NewSocket : Constructor
func NewSocket(conn *net.Conn) *Socket {
	s := Socket{
		events:      make(map[string]func(messages.Payload)),
		id:          generateSocketID(),
		conn:        conn,
		initialized: false,
	}
	go s.listen()
	return &s
}

func (s *Socket) initialize(payload messages.Payload) {
	val, ok := payload["RoomName"]
	if ok {
		s.Join(val)
	}
	s.initialized = true
}

func (s *Socket) matchMaking(payload messages.Payload) {
	s.JoinMatchMaking()
}

func (s *Socket) listen() {
	reader := bufio.NewReader(*s.conn)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		msg, err := messages.NewMessage(text)
		if err != nil {
			continue
		}
		if msg.Name == "Init" {
			s.initialize(msg.Body)
			continue
		}
		if msg.Name == "MatchMaking" {
			s.matchMaking(msg.Body)
			continue
		}
		f, ok := s.events[msg.Name]
		if ok && s.IsOk() {
			f(msg.Body)
		} else {
			fmt.Println("Event name doesn't found")
		}
	}
}

func generateSocketID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
