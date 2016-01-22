package tank

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var refreshModifier float32 = 1

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	clients   map[int]*Client
	bullets   []*Bullet
	explosion Explosion
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
	score     Scores
	mapa      *mapa
}

// Create new chat server.
func NewServer(pattern string, mod float32) *Server {
	var bullets []*Bullet
	messages := []*Message{}
	explosion := Explosion{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)
	var score Scores
	score.client = make(map[int]int)
	refreshModifier = mod
	m := &mapa{}

	s := &Server{
		pattern,
		messages,
		clients,
		bullets,
		explosion,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
		score,
		m,
	}

	s.mapa.setMap()
	return s
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

// func (s *Server) Done() {
// 	s.doneCh <- true
// }

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	x := s.BuildAnswer(c.id, true)
	c.Write(&x)
}

func (s *Server) sendAll() {
	s.calcAll()
	for _, c := range s.clients {
		m := s.BuildAnswer(c.id, false)
		c.Write(&m)
	}
	s.scoreRead()
}

func (s *Server) RunInterval(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			s.sendAll()
			//		default:
			//			log.Println("wtf stop")
			//			ticker.Stop()
			//			return
		}
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s, len(s.clients))
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			log.Println("asdas")
			return
		}
	}
}
