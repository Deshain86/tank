package engine

import (
	"fmt"
	"net"
)

const channelBufSize int = 100

var maxId int = 0
var fullLife int = 100
var defaultDirection int = 0
var prev string

var firstPosition [][]float32 = [][]float32{
	[]float32{50, 50},
	[]float32{canvasSizeX - 100, canvasSizeY - 100},
	[]float32{50, canvasSizeY - 100},
	[]float32{canvasSizeX - 100, 50}}

type Client struct {
	id            int
	nick          string
	RemoteAddr    *net.UDPAddr
	RemoteAddrStr string
	server        *Server
	ch            chan *string
	doneCh        chan bool
	PositionX     float32
	PositionY     float32
	Life          int
	Direction     int
	Speed         float32
	Moving        bool
	Fire          bool
	LastFire      int
	StartPosX     float32
	StartPosY     float32
}

// Create new chat client.
func (server *Server) NewClient(remoteAddr *net.UDPAddr, nick string) *Client {
	if remoteAddr == nil {
		panic("remoteAddr cannot be nil")
	}

	maxId = len(server.clients) + 1
	ch := make(chan *string, channelBufSize)
	doneCh := make(chan bool)
	position := firstPosition[maxId%4]

	return &Client{
		maxId,
		nick,
		remoteAddr,
		remoteAddr.String(),
		server,
		ch,
		doneCh,
		float32(position[0]),
		float32(position[1]),
		fullLife,
		defaultDirection,
		defaultTankSpeed,
		false,
		false,
		0,
		float32(position[0]),
		float32(position[1])}
}

func (c *Client) GetId() int {
	return c.id
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
