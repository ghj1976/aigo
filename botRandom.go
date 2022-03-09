package aigo

import (
	"log"
	"math/rand"
)

type IAgent interface {
	SelectMove(gs *GameState) Move
}

// 随机下棋机器人
type RandomBot struct {
	IAgent
}

func (bot RandomBot) SelectMove(gs *GameState) Move {
	candidates := []Point{} // 候选可以下棋的点

	for r := uint16(1); r <= uint16(gs.BoardPosition.Height); r++ {

		for c := uint16(1); c <= uint16(gs.BoardPosition.Width); c++ {
			candidate := Point{Row: r, Col: c}
			if gs.IsValidMove(NewPlay(candidate)) {
				// log.Println("pass IsValidMove")
				if !gs.BoardPosition.IsPointAnEye(candidate, gs.PlayerTurn) {
					// log.Println("pass IsPointAnEye")
					candidates = append(candidates, candidate)
				}
			}
		}
	}
	if len(candidates) <= 0 {
		log.Printf("候选清单:%d", len(candidates))
		return NewPass()
	}

	return NewPlay(candidates[rand.Int31n(int32(len(candidates)))])
}
