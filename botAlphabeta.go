package aigo

import (
	"log"
	"math/rand"
)

// 使用评估函数+ αβ剪枝算法 的机器人
type AlphaBetaAgent struct {
	MaxDepth int                     // 最大搜索深度
	EvalFn   func(gs *GameState) int // 棋局评估函数
}

//  将函数作为参数 的例子 https://www.kancloud.cn/kancloud/the-way-to-go/72479
func NewAlphaBetaAgent(depth int, evalFn func(gs *GameState) int) *AlphaBetaAgent {
	bot := &AlphaBetaAgent{}
	bot.MaxDepth = depth
	bot.EvalFn = evalFn
	return bot
}

// 机器人选择下一步如何走
func (bot *AlphaBetaAgent) SelectMove(gs *GameState) Move {
	best_moves := []Move{}
	best_score := MIN_SCORE
	best_black := MIN_SCORE
	best_white := MIN_SCORE
	for _, possible_move := range gs.LegalMoves() {
		next_state, err := gs.ApplyMove(possible_move)
		if err != nil {
			log.Fatalln(err)
		}
		opponent_best_outcome := next_state.AlphaBetaResult(bot.MaxDepth, best_black, best_white, bot.EvalFn)
		our_best_outcome := -1 * opponent_best_outcome

		if len(best_moves) <= 0 || our_best_outcome > best_score {
			best_moves = []Move{possible_move} // 清空原先已有的，以算出来最佳覆盖
			best_score = our_best_outcome

			// 下面是 alpha-beta 剪枝算法的关键
			if gs.PlayerTurn == Black {
				best_black = best_score
			} else if gs.PlayerTurn == White {
				best_white = best_score
			}

		} else if our_best_outcome == best_score {
			best_moves = append(best_moves, possible_move)
		}

	}
	// 找出棋局评估函数评估价值最大的一步
	return best_moves[rand.Int31n(int32(len(best_moves)))]
}

// 通过αβ剪枝算法，棋局评估函数 ，找最佳走法
// 多了2个输入参数， 目前的 best_black, best_white
func (gs *GameState) AlphaBetaResult(max_depth, best_black, best_white int, evalFn func(gs *GameState) int) int {
	if gs.IsOver() {
		if gs.Winner() == gs.PlayerTurn {
			return MAX_SCORE
		} else {
			return MIN_SCORE
		}
	}
	if max_depth == 0 { // 超过最大递归深度后的采用棋局评估函数
		return evalFn(gs)
	}

	best_so_far := MIN_SCORE

	for _, candidate_move := range gs.LegalMoves() {
		next_state, err := gs.ApplyMove(candidate_move)
		if err != nil {
			log.Panicln(err)
		}
		// 递归自身
		opponent_best_result := next_state.AlphaBetaResult(max_depth-1, best_black, best_white, evalFn)

		our_result := -1 * opponent_best_result
		if our_result > best_so_far {
			best_so_far = our_result
		}

		// 下面是 alpha-beta 剪枝算法的关键， 跳过一些更差的判断
		if gs.PlayerTurn == White {
			if best_so_far > best_white {
				best_white = best_so_far
			}
			outcome_for_blank := -1 * best_so_far
			if outcome_for_blank < best_black {
				return best_so_far // 剪枝 比目前黑棋更差的跳过
			}

		} else if gs.PlayerTurn == Black {
			if best_so_far > best_black {
				best_black = best_so_far
			}
			outcome_for_white := -1 * best_so_far
			if outcome_for_white < best_white {
				return best_so_far //  剪枝 比目前白棋更差的跳过
			}
		}

	}

	return best_so_far
}
