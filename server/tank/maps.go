package tank

import (
	"log"
)

type mapa struct {
	ground     [][]int
	speedPoint [][]int
}

func (s *Server) getSpeedPosition(x, y float32) float32 {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
			log.Print(x, y)
		}
	}()
	// log.Print(x, y)
	return float32(s.mapa.speedPoint[int(x)][int(y)]) / 10
}

func (s *Server) getCollision(x, y float32) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
			log.Print(x, y)
		}
	}()
	// log.Print(x, y)
	return s.mapa.speedPoint[int(x)][int(y)] == 0
}

func (s *mapa) drawMap() [][]int {
	return s.ground
}

func (s *mapa) setMap() {
	s.ground = mapa1

	res := make([][]int, 800)
	for x := 0; x < len(res); x++ {
		res[x] = make([]int, 800)
	}

	for k1, v1 := range s.ground {
		for k2, v2 := range v1 {
			speed := speedGround[v2]
			for j := 0; j < 50; j++ {
				for i := 0; i < 50; i++ {
					res[(k2*50)+j][(k1*50)+i] = speed
				}
			}
		}
	}

	// // woda
	// for x := 280; x < 380; x++ {
	// 	for y := 580; y < 800; y++ {
	// 		res[x][y] = 3
	// 	}
	// }
	// // lod
	// for x := 280; x < 380; x++ {
	// 	for y := 280; y < 480; y++ {
	// 		res[x][y] = 30
	// 	}
	// }
	// // lod
	// for x := 380; x < 480; x++ {
	// 	for y := 280; y < 380; y++ {
	// 		res[x][y] = 30
	// 	}
	// }
	// // drzewo
	// for x := 167; x < 396; x++ {
	// 	for y := 67; y < 200; y++ {
	// 		res[x][y] = 0
	// 	}
	// }
	// // drzewo
	// for x := 267; x < 396; x++ {
	// 	for y := 0; y < 67; y++ {
	// 		res[x][y] = 0
	// 	}
	// }
	s.speedPoint = res
}

var speedGround []int = []int{10, 0, 10, 30, 3}

var mapa1 [][]int = [][]int{
	[]int{2, 2, 2, 0, 0, 2, 0, 1, 2, 0, 0, 2, 0, 2, 0, 2},
	[]int{0, 2, 0, 0, 2, 0, 2, 1, 2, 2, 0, 2, 0, 2, 2, 0},
	[]int{0, 2, 0, 0, 2, 2, 2, 1, 2, 0, 2, 2, 0, 2, 0, 2},
	[]int{0, 2, 0, 0, 2, 0, 2, 1, 2, 0, 0, 2, 0, 2, 0, 2},
	[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
	[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
	[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
	[]int{3, 0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 3},
	[]int{3, 0, 0, 0, 0, 4, 4, 0, 0, 4, 4, 0, 0, 0, 0, 3},
	[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
	[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
	[]int{3, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 3},
	[]int{0, 2, 0, 0, 0, 2, 0, 1, 2, 0, 0, 2, 0, 2, 2, 2},
	[]int{2, 0, 0, 0, 2, 0, 2, 1, 2, 2, 2, 2, 0, 2, 0, 0},
	[]int{2, 2, 2, 0, 2, 2, 2, 1, 2, 0, 0, 2, 0, 2, 2, 2},
	[]int{0, 2, 2, 0, 2, 0, 2, 1, 2, 0, 0, 2, 0, 2, 2, 2},
}
