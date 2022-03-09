package aigo

import (
	"math/rand"
	"time"
)

// 加速的随机下棋机器人
type FastRandomBot struct {
	point_cache []Point // 缓存的棋盘点集合，每次都需要把它打散下
}

func NewFastRandomBot() *FastRandomBot {
	bot := &FastRandomBot{}
	bot.point_cache = []Point{}
	return bot
}

// 随机打乱数组
// https://golangnote.com/topic/260.html
func randShuffle(slice []Point) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func (bot *FastRandomBot) SelectMove(gs *GameState) Move {

	if len(bot.point_cache) <= 0 {
		for r := 1; r <= int(gs.BoardPosition.Height); r++ {
			for c := 1; c <= int(gs.BoardPosition.Width); c++ {
				bot.point_cache = append(bot.point_cache, Point{Row: uint16(r), Col: uint16(c)})
			}
		}
	}

	randShuffle(bot.point_cache)
	for i := 0; i < len(bot.point_cache); i++ {
		candidate := bot.point_cache[i]

		if gs.IsValidMove(NewPlay(candidate)) {
			if !gs.BoardPosition.IsPointAnEye(candidate, gs.PlayerTurn) {
				return NewPlay(candidate)
			}
		}
	}
	return NewPass() // 跳过
}
