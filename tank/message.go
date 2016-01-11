package tank

type Message struct {
	Author    string `json:"author"`
	Body      string `json:"body"`
	PositionX int    `json:"positionX"`
	PositionY int    `json:"positionY"`
}

type Answer struct {
	Users []User `json:"users"`
	Body  string `json:"body"`
}

type User struct {
	Id        int    `json:"id"`
	Color     string `json:"color"`
	PositionX int    `json:"posX"`
	PositionY int    `json:"posY"`
}

// func (self *Message) ParseResponse(clientId int) {
func (self *Server) ParseResponse(msg *Message, clientId int) {
	tmp := self.clients[clientId] // users[clientId]
	switch msg.Body {
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
	case "right2":
		tmp.Direction = 90
		tmp.Moving = false
	case "left2":
		tmp.Direction = 270
		tmp.Moving = false
	case "down2":
		tmp.Direction = 180
		tmp.Moving = false
	case "up2":
		tmp.Direction = 0
		tmp.Moving = false
	}
	self.clients[clientId] = tmp
}

func (self *Server) BuildAnswer(clientId int) Answer {
	var ans Answer
	for _, user := range self.clients {
		var u User
		u.Id = user.id
		if clientId == user.id {
			u.Color = "b"
		} else {
			u.Color = "r"
		}

		if user.Moving {
			switch user.Direction {
			case 0:
				user.PositionY--
			case 90:
				user.PositionX++
			case 180:
				user.PositionY++
			case 270:
				user.PositionX--
			}
		}

		u.PositionX = user.PositionX
		u.PositionY = user.PositionY
		ans.Users = append(ans.Users, u)
	}
	return ans
}
