package tank

const bulletSpeed float32 = 6

const bulletWidthHalf float32 = 3
const bulletHeightHalf float32 = 6.5

type Bullet struct {
	x         float32
	y         float32
	direction int
	ownerId   int
}

func (s *Server) checkBulletsOnMap(mapSizeX, mapSizeY float32, refreshTime float32) {
	var bSpeed = bulletSpeed * refreshTime
	var newList []*Bullet
forLoop:
	for _, b := range s.bullets {
		switch b.direction {
		case 0:
			b.y -= bSpeed
			if b.y < 0 {
				continue forLoop
			}
		case 90:
			b.x += bSpeed
			if b.x > mapSizeX {
				continue forLoop
			}
		case 180:
			b.y += bSpeed
			if b.y > mapSizeY {
				continue forLoop
			}
		case 270:
			b.x -= bSpeed
			if b.x < 0 {
				continue forLoop
			}
		}
		newList = append(newList, b)
	}
	s.bullets = newList
}

func (s *Server) checkHitBullet(clientId int, tankX1, tankY1, tankX2, tankY2 float32) bool {
	for k, b := range s.bullets {
		if b.ownerId != clientId {
			if (tankX2 > b.x && tankX1 < b.x) && (tankY2 > b.y && tankY1 < b.y) {
				var tmpList []*Bullet
				for k2, b2 := range s.bullets {
					if k != k2 {
						tmpList = append(tmpList, b2)
					}
				}
				s.bullets = tmpList
				return true
			}
		}
	}
	return false
}
