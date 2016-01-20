package tank

const tankWidth float32 = 37.5
const tankWidthHalf float32 = 18.75
const tankHeight float32 = 35
const tankHeightHalf float32 = 17.5

func (s *Server) checkHitTank(c *Client) bool {
	if s.checkHitBullet(c.id, c.PositionX, c.PositionY, c.PositionX+tankWidth, c.PositionY+tankHeight) {
		return true
	}
	return false
}
