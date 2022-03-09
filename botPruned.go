package aigo

import (
	"log"
	"math/rand"
)

const (
	MAX_SCORE = 999999
	MIN_SCORE = -999999
)

// 使用评估函数的剪枝算法机器人
type DepthPrunedAgent struct {
	MaxDepth int                     // 最大搜索深度
	EvalFn   func(gs *GameState) int // 棋局评估函数
}

//  将函数作为参数 的例子 https://www.kancloud.cn/kancloud/the-way-to-go/72479
func NewDepthPrunedAgent(depth int, evalFn func(gs *GameState) int) *DepthPrunedAgent {
	bot := &DepthPrunedAgent{}
	bot.MaxDepth = depth
	bot.EvalFn = evalFn
	return bot
}

// 机器人选择一步
func (bot *DepthPrunedAgent) SelectMove(gs *GameState) Move {
	best_moves := []Move{}
	best_score := MIN_SCORE
	for _, possible_move := range gs.LegalMoves() {
		next_state, err := gs.ApplyMove(possible_move)
		if err != nil {
			log.Fatalln(err)
		}
		opponent_best_outcome := next_state.BestResult(bot.MaxDepth, bot.EvalFn)
		our_best_outcome := -1 * opponent_best_outcome

		if len(best_moves) <= 0 || our_best_outcome > best_score {
			best_moves = []Move{possible_move} // 清空原先已有的，以算出来最佳覆盖
			best_score = our_best_outcome
		} else if our_best_outcome == best_score {
			best_moves = append(best_moves, possible_move)
		}

	}
	// 找出棋局评估函数评估价值最大的一步
	return best_moves[rand.Int31n(int32(len(best_moves)))]
}

// 通过剪枝算法，棋局评估函数 ，找最佳走法
func (gs *GameState) BestResult(max_depth int, evalFn func(gs *GameState) int) int {
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

		opponent_best_result := next_state.BestResult(max_depth-1, evalFn)

		our_result := -1 * opponent_best_result
		if our_result > best_so_far {
			best_so_far = our_result
		}

	}

	return best_so_far
}
