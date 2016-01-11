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
		tmp.PositionX++
	case "left":
		tmp.PositionX--
	case "down":
		tmp.PositionY++
	case "up":
		tmp.PositionY--
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
		u.PositionX = user.PositionX
		u.PositionY = user.PositionY
		ans.Users = append(ans.Users, u)
	}
	return ans
}
