package tank

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxId int = 0
var fullLife int = 100
var defaultDirection int = 0

type Client struct {
	id        int
	ws        *websocket.Conn
	server    *Server
	ch        chan *Answer
	doneCh    chan bool
	PositionX int
	PositionY int
	Life      int
	Direction int
	Moving    bool
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
	ch := make(chan *Answer, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, ws, server, ch, doneCh, 10, 10, fullLife, defaultDirection, false}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(ans *Answer) {
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
			log.Println("Send:", msg)
			websocket.Message.Send(c.ws, buildAnswer(msg))

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

func buildAnswer(msg *Answer) string {
	var result string
	for _, u := range msg.Users {
		result += fmt.Sprintf("T;%d;%s;%d;%d;%d;%d;\n",
			u.Id, u.Color, u.PositionX, u.PositionY, 0, 100)
	}
	return result
}

/*
Odpowiedz format
tank
obiekt;id;color;pozycjaX;pozycjaY;obrot;zycie(hp);
T;1;R;10;10;360;50;

kolor R G B K

*/
