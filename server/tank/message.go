package tank

import (
	"bytes"
	"fmt"
	"strconv"
)

type Message struct {
	// Author    string
	Body string
	// PositionX int
	// PositionY int
}

type Answer struct {
	Users   []User
	Bullets []*Bullet
}

type User struct {
	Id        int
	Color     string
	PositionX float32
	PositionY float32
	Direction int
}

func (self *Server) ParseResponse(msg *string, clientId int) {
	tmp := self.clients[clientId] // users[clientId]
	switch *msg {
	case "fire":
		tmp.Fire = true
	case "fire2":
		tmp.Fire = false
	case "right":
		tmp.Direction = 90
		tmp.Moving = true
		tmp.Speed = defaultTankSpeed
	case "left":
		tmp.Direction = 270
		tmp.Moving = true
		tmp.Speed = defaultTankSpeed
	case "down":
		tmp.Direction = 180
		tmp.Moving = true
		tmp.Speed = defaultTankSpeed
	case "up":
		tmp.Direction = 0
		tmp.Moving = true
		tmp.Speed = defaultTankSpeed
	case "right2", "left2", "down2", "up2":
		tmp.Moving = false
		tmp.Speed = 0
	}
	self.clients[clientId] = tmp
}

func (self *Server) BuildAnswer(clientId int, firstAnswer bool) string {
	var result bytes.Buffer
	if firstAnswer {
		x := self.mapa.drawMap()
		for _, v := range x {
			result.WriteString("M;")
			for _, v2 := range v {
				result.WriteString(strconv.Itoa(v2) + ";")
			}
			result.WriteString("\n")
		}
	}
	for _, u := range self.bullets {
		result.WriteString(fmt.Sprintf("B;%.0f;%.0f;%d;\n",
			u.x, u.y, u.direction))
	}
	for _, user := range self.clients {
		color := "r"
		if clientId == user.id {
			color = "b"
		}

		result.WriteString(fmt.Sprintf("T;%d;%s;%.0f;%.0f;%.0f;%d;%d;%d;\n",
			user.id, color, user.PositionX, user.PositionY, user.Speed, user.Direction, user.Direction, 100))
	}
	if self.score.change {
		for id, point := range self.score.client {
			result.WriteString(fmt.Sprintf("S;%d;%d;\n", id, point))
		}
	}
	return result.String()
}

/*
Odpowiedz format
tank
obiekt;id;color;pozycjaX;pozycjaY;obrot;obrot_lufy;zycie(hp);
T;1;R;10;10;0;0;50;

kolor R G B K

*/
