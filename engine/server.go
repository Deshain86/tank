package engine

import (
	"log"
	"net"
	"strconv"
)

var refreshModifier float32 = 1

// Chat server.
type Server struct {
	conn      *net.UDPConn
	messages  []*Message
	clients   map[string]*Client
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
func NewServer(conn *net.UDPConn) *Server {
	var bullets []*Bullet
	messages := []*Message{}
	explosion := Explosion{}
	clients := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)
	var score Scores
	score.client = make(map[int]int)
	//	refreshModifier = mod
	m := &mapa{}

	s := &Server{
		conn,
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
	s.sendResponse(c.RemoteAddr, strconv.Itoa(c.GetId()))
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	x := s.BuildAnswer(c.id, true)
	c.Write(&x)
}

func (s *Server) SendAll() {
	s.calcAll()
	for _, c := range s.clients {
		m := s.BuildAnswer(c.id, false)
		s.sendResponse(c.RemoteAddr, m)
		log.Print("msg: ", m)
	}
	s.scoreRead()
	s.explosionRead()
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.RemoteAddrStr] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.RemoteAddrStr)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			log.Println("asdas")
			return
		}
	}
}

func (s *Server) sendResponse(addr *net.UDPAddr, msg string) {
	_, err := s.conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		log.Printf("Couldn't send response %v", err)
	}
}
