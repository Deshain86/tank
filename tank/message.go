package tank

import (
	"fmt"
)

type Message struct {
	Author    string
	Body      string
	PositionX int
	PositionY int
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

// func (self *Message) ParseResponse(clientId int) {
func (self *Server) ParseResponse(msg *Message, clientId int) {
	tmp := self.clients[clientId] // users[clientId]
	switch msg.Body {
	case "space":
		tmp.Fire = true
	case "space2":
		tmp.Fire = false
	case "right":
		tmp.Direction = 90
		tmp.Moving = true
	case "left":
		tmp.Direction = 270
		tmp.Moving = true
	case "down":
		tmp.Direction = 180
		tmp.Moving = true
	case "up":
		tmp.Direction = 0
		tmp.Moving = true
	case "right2", "left2", "down2", "up2":
		tmp.Moving = false
	}
	self.clients[clientId] = tmp
}

func (self *Server) BuildAnswer(clientId int) string {
	var result string
	for _, u := range self.bullets {
		result += fmt.Sprintf("B;%d;%d;\n",
			u.x, u.y)
	}
	for _, user := range self.clients {
		color := "r"
		if clientId == user.id {
			color = "b"
		}

		result += fmt.Sprintf("T;%d;%s;%f;%f;%d;%d;%d;\n",
			user.id, color, user.PositionX, user.PositionY, user.Direction, user.Direction, 100)
	}

	return result
}

/*
Odpowiedz format
tank
obiekt;id;color;pozycjaX;pozycjaY;obrot;obrot_lufy;zycie(hp);
T;1;R;10;10;0;0;50;

kolor R G B K

*/
