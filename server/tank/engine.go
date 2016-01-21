package tank

const canvasSizeX float32 = 800
const canvasSizeY float32 = 800

func (s *Server) calcAll() {
	s.checkBulletsOnMap(canvasSizeX, canvasSizeY, refreshModifier)

forLoop:
	for _, c := range s.clients {
		hit, hitClientId := s.checkHitTank(c)
		if hit {
			c.PositionX = c.StartPosX
			c.PositionY = c.StartPosY
			s.scoreAdd(hitClientId)
			m := s.BuildAnswer(c.id)
			c.Write(&m)
			continue forLoop
		}

		if c.Fire {
			if c.LastFire == 0 {
				c.LastFire = 20 * int(refreshModifier)
				s.bullets = append(s.bullets,
					&Bullet{
						ownerId:   c.id,
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
	}
}
