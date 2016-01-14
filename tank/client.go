package tank

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize int = 100

var maxId int = 0
var fullLife int = 100
var defaultDirection int = 0
var prev string

type Client struct {
	id        int
	ws        *websocket.Conn
	server    *Server
	ch        chan *string
	doneCh    chan bool
	PositionX float32
	PositionY float32
	Life      int
	Direction int
	Speed     float32
	Moving    bool
	Fire      bool
	LastFire  int
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan *string, channelBufSize)
	doneCh := make(chan bool)

	return &Client{
		maxId,
		ws,
		server,
		ch,
		doneCh,
		float32(10),
		float32(10),
		fullLife,
		defaultDirection,
		defaultTankSpeed,
		false,
		false,
		0}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(ans *string) {
	select {
	case c.ch <- ans:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			websocket.Message.Send(c.ws, *msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				log.Print("MSG ", msg)
				c.server.SendAll(&msg, c.id)
			}
		}
	}
}
