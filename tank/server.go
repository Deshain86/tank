package tank

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

const bulletSpeed float32 = 4
const canvasSizeX float32 = 800
const canvasSizeY float32 = 800
const tankWidth float32 = 37
const tankHeight float32 = 35
const defaultTankSpeed float32 = 2

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	clients   map[int]*Client
	bullets   []*Bullet
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// var users map[int]Positions = make(map[int]Positions)

// type Positions struct {
// 	x int
// 	y int
// }

type Bullet struct {
	x         float32
	y         float32
	direction int
}

// Create new chat server.
func NewServer(pattern string) *Server {
	var bullets []*Bullet
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		clients,
		bullets,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *Message, clientId int) {
	s.ParseResponse(msg, clientId)
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	x := s.BuildAnswer(c.id)
	c.Write(&x)
}

func (s *Server) sendAll() {
	// log.Print("CNT ", len(s.bullets), len(s.clients))
	for k, b := range s.bullets {
		switch b.direction {
		case 0:
			b.y -= bulletSpeed
			if b.y > canvasSizeY || b.y < 0 {
				s.bullets = append(s.bullets[:k], s.bullets[k+1:]...)
			}
		case 90:
			b.x += bulletSpeed
			if b.x > canvasSizeX || b.x < 0 {
				s.bullets = append(s.bullets[:k], s.bullets[k+1:]...)
			}
		case 180:
			b.y += bulletSpeed
			if b.y > canvasSizeY || b.y < 0 {
				s.bullets = append(s.bullets[:k], s.bullets[k+1:]...)
			}
		case 270:
			b.x -= bulletSpeed
			if b.x > canvasSizeX || b.x < 0 {
				s.bullets = append(s.bullets[:k], s.bullets[k+1:]...)
			}
		}
	}

	for _, c := range s.clients {
		if c.Fire {
			if c.LastFire == 0 {
				c.LastFire = 10
				s.bullets = append(s.bullets, &Bullet{x: c.PositionX, y: c.PositionY, direction: c.Direction})
				// log.Print("FIRE ", len(s.bullets))
			}
			// log.Print("TRY FIRE ")
			c.LastFire--
		}
		if c.Moving {
			switch c.Direction {
			case 0:
				c.PositionY = c.PositionY - c.Speed
				if c.PositionY <= 0 {
					c.PositionY = 0
				}
			case 90:
				c.PositionX = c.PositionX + c.Speed
				log.Println(c.PositionX, tankHeight, canvasSizeX)
				if c.PositionX+tankHeight >= canvasSizeX {
					c.PositionX = canvasSizeX - tankHeight
				}
			case 180:
				c.PositionY = c.PositionY + c.Speed
				if c.PositionY+tankHeight >= canvasSizeY {
					c.PositionY = canvasSizeY - tankHeight
				}
			case 270:
				c.PositionX = c.PositionX - c.Speed
				if c.PositionX <= 0 {
					c.PositionX = 0
				}
			}
		}
		m := s.BuildAnswer(c.id)
		c.Write(&m)
	}
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

		client := NewClient(ws, s)
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
