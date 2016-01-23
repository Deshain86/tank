package tank

import (
	"log"
)

type mapa struct {
	speedMap [][]int
}

func (s *Server) getSpeedPosition(x, y float32) float32 {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
			log.Print(x, y)
		}
	}()
	// log.Print(x, y)
	return float32(s.mapa.speedMap[int(x)][int(y)]) / 10
}

func (s *Server) getCollision(x, y float32) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
			log.Print(x, y)
		}
	}()
	// log.Print(x, y)
	return s.mapa.speedMap[int(x)][int(y)] == 0
}

func (s *Server) setMap() {
	res := make([][]int, 800)
	for x := 0; x < len(res); x++ {
		res[x] = make([]int, 800)
		for y := 0; y < len(res[x]); y++ {
			res[x][y] = 10
		}
	}

	// woda
	for x := 280; x < 380; x++ {
		for y := 580; y < 800; y++ {
			res[x][y] = 3
		}
	}
	// lod
	for x := 280; x < 380; x++ {
		for y := 280; y < 480; y++ {
			res[x][y] = 30
		}
	}
	// lod
	for x := 380; x < 480; x++ {
		for y := 280; y < 380; y++ {
			res[x][y] = 30
		}
	}
	// drzewo
	for x := 167; x < 396; x++ {
		for y := 67; y < 200; y++ {
			res[x][y] = 0
		}
	}
	// drzewo
	for x := 267; x < 396; x++ {
		for y := 0; y < 67; y++ {
			res[x][y] = 0
		}
	}
	s.mapa.speedMap = res
}
