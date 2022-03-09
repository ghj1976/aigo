package ttt

import (
	"log"
	"math/rand"
)

type MinimaxAgent struct {
}

func (agent *MinimaxAgent) SelectMove(game_state GameState) Move {

	// 注意：前面这些判断， 图书给的例子中并不包含， 缺少这部分判断。
	// 原先逻辑有问题：未来可以赢，并不代表最近的一步走错就会输。

	// 先检查是否可以在下一步直接获胜。如果可以，就这样行动。
	move := game_state.Find_Winning_Move(game_state.PlayerTurn)
	if move != nil {
		return *move
	}

	// 如果不可以，再看看对手能否在下一步获胜。如果能，就尝试阻止它。
	moveArr := game_state.Eliminate_Losing_Move(game_state.PlayerTurn)
	lenMA := len(moveArr)
	if lenMA > 0 {
		return moveArr[rand.Intn(lenMA)]
	}

	// 如果对手不能获胜，再看看能否通过第二步棋取胜。如果能，就按照这两步棋来落子。
	move2 := game_state.Find_Two_Step_Win(game_state.PlayerTurn)
	if move2 != nil {
		return *move2
	}
	// 如果不能，再看看对手的第二步棋是否能获胜。

	// 递归找更多步
	winning_moves := []Move{}
	draw_moves := []Move{}
	losing_moves := []Move{}
	for _, possible_move := range game_state.GetAllLegalMoves() {
		next_state := game_state.ApplyMove(possible_move)

		opponent_best_result := next_state.BestResult()
		if opponent_best_result == Win {
			winning_moves = append(winning_moves, possible_move)
		} else if opponent_best_result == Draw {
			draw_moves = append(draw_moves, possible_move)
		} else {
			losing_moves = append(losing_moves, possible_move)
		}

	}

	if lenw := len(winning_moves); lenw > 0 {
		log.Printf("Win:%v", winning_moves)
		return winning_moves[rand.Intn(lenw)]
	}

	if lend := len(draw_moves); lend > 0 {
		log.Printf("draw:%v", draw_moves)
		return draw_moves[rand.Intn(lend)]
	}

	lenl := len(losing_moves)
	log.Printf("losing:%v", losing_moves)
	return losing_moves[rand.Intn(lenl)]
}
