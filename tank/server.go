package tank

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

const bulletSpeed float32 = 6
const canvasSizeX float32 = 800
const canvasSizeY float32 = 800

const tankWidth float32 = 37.5
const tankWidthHalf float32 = 18.75
const tankHeight float32 = 35
const tankHeightHalf float32 = 17.5

const bulletWidthHalf float32 = 3
const bulletHeightHalf float32 = 6.5

const defaultTankSpeed float32 = 2

var refreshModifier float32 = 1

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

type Bullet struct {
	x         float32
	y         float32
	direction int
}

// Create new chat server.
func NewServer(pattern string, mod float32) *Server {
	var bullets []*Bullet
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)
	refreshModifier = mod

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
	var bSpeed = bulletSpeed * refreshModifier
	var deleteList []int
	for k, b := range s.bullets {
		switch b.direction {
		case 0:
			b.y -= bSpeed
			if b.y < 0 {
				deleteList = append(deleteList, k)
			}
		case 90:
			b.x += bSpeed
			if b.x > canvasSizeX {
				deleteList = append(deleteList, k)
			}
		case 180:
			b.y += bSpeed
			if b.y > canvasSizeY {
				deleteList = append(deleteList, k)
			}
		case 270:
			b.x -= bSpeed
			if b.x < 0 {
				deleteList = append(deleteList, k)
			}
		}
	}
	// log.Print(deleteList)
	if len(deleteList) > 0 {
		var tmp []*Bullet
	forLoop:
		for k, b := range s.bullets {
			for _, del := range deleteList {
				if del == k {
					continue forLoop
				}
			}
			tmp = append(tmp, b)
		}
		s.bullets = tmp
		// log.Print(len(s.bullets), len(tmp), len(deleteList))
	}

	for _, c := range s.clients {
		if c.Fire {
			if c.LastFire == 0 {
				c.LastFire = 20 * int(refreshModifier)
				s.bullets = append(s.bullets,
					&Bullet{
						x:         c.PositionX + tankWidthHalf - bulletWidthHalf,
						y:         c.PositionY + tankHeightHalf - bulletHeightHalf,
						direction: c.Direction})
			}
		}
		if c.LastFire > 0 {
			c.LastFire--
		}

		var speed = c.Speed * refreshModifier
		if c.Moving {
			switch c.Direction {
			case 0:
				c.PositionY = c.PositionY - speed
				if c.PositionY <= 0 {
					c.PositionY = 0
				}
			case 90:
				c.PositionX = c.PositionX + speed
				if c.PositionX+tankHeight >= canvasSizeX {
					c.PositionX = canvasSizeX - tankHeight
				}
			case 180:
				c.PositionY = c.PositionY + speed
				if c.PositionY+tankHeight >= canvasSizeY {
					c.PositionY = canvasSizeY - tankHeight
				}
			case 270:
				c.PositionX = c.PositionX - speed
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
