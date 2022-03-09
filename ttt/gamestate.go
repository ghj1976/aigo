package ttt

import "log"

type GameState struct {
	BoardPosition *Board // 棋子的布局
	PlayerTurn    Player // 当前是那个玩家的回合
	LastMove      *Move  // 上一步动作
}

func NewGameState() *GameState {
	gs := &GameState{}
	gs.BoardPosition = NewBoard()
	gs.PlayerTurn = X
	gs.LastMove = nil

	return gs
}

// 获得所有的合法移动
func (gs *GameState) GetAllLegalMoves() []Move {
	moves := []Move{}
	for row := 1; row <= BOARD_SIZE; row++ {
		for col := 1; col <= BOARD_SIZE; col++ {
			move := Move{P: Point{uint16(row), uint16(col)}}
			if gs.IsValidMove(move) {
				moves = append(moves, move)
			}
		}
	}

	return moves
}

// 没结束，并且有位置，才可以继续移动
func (gs *GameState) IsValidMove(m Move) bool {
	if gs.BoardPosition.Get(m.P) == None && !gs.IsOver() {
		return true
	}
	return false
}

func (gs *GameState) IsOver() bool {
	if gs.Winner() != None {
		return true
	}

	// 没有地方可以下棋了，也是结束了
	for row := 1; row <= BOARD_SIZE; row++ {
		for col := 1; col <= BOARD_SIZE; col++ {
			if gs.BoardPosition.Get(Point{uint16(row), uint16(col)}) == None {
				return false
			}
		}
	}
	return true
}

// 是否有人赢了？
func (gs *GameState) Winner() Player {

	if gs.BoardPosition.Has_3_in_a_row(X) {
		return X
	}
	if gs.BoardPosition.Has_3_in_a_row(O) {
		return O
	}
	return None
}

func (gs *GameState) ApplyMove(move Move) *GameState {
	next_board := gs.BoardPosition.DeepCopy()
	next_board.Place(gs.PlayerTurn, move.P)

	log.Printf("%v:%v", gs.PlayerTurn, move.P)
	return &GameState{BoardPosition: next_board,
		PlayerTurn: gs.PlayerTurn.Other(),
		LastMove:   &move,
	}
}

func (gs *GameState) BestResult() GameResult {
	if gs.IsOver() {
		w := gs.Winner()
		log.Printf("结果:%s", w.String())
		switch w {
		case gs.PlayerTurn:
			return Win
		case None:
			return Draw
		default:
			return Loss
		}
	}

	// 递归找最好的步骤
	best_result_so_far := Loss
	for _, candidate_move := range gs.GetAllLegalMoves() {
		next_state := gs.ApplyMove(candidate_move)               // 看看走出这一步，棋局会如何变化？
		opponent_best_result := next_state.BestResult()          // 找到对方的最佳动作
		our_result := opponent_best_result.reverse_game_result() // 无论对方想要什么？我们想要的都是它的反面
		if our_result > best_result_so_far {                     // 看看当前的结果是否比之前得到的结果更好？
			best_result_so_far = our_result
		}
	}

	return best_result_so_far
}

// 寻找一步就可以赢的步骤， 返回为nil则没有一步就赢的
// 如果立即可以赢的有多个位置，随便返回一个
func (gs *GameState) Find_Winning_Move(nextPlayer Player) *Move {
	for _, candidate_move := range gs.GetAllLegalMoves() {
		next_state := gs.ApplyMove(candidate_move) // 看看走出这一步，棋局会如何变化？
		if next_state.IsOver() && next_state.Winner() == nextPlayer {
			return &candidate_move
		}
	}
	return nil
}

// 寻找避免让对手直接获胜的
func (gs *GameState) Eliminate_Losing_Move(nextPlayer Player) []Move {
	possible_moveArr := []Move{}
	opponent := nextPlayer.Other()
	for _, candidate_move := range gs.GetAllLegalMoves() {
		next_state := gs.ApplyMove(candidate_move) // 看看走出这一步，棋局会如何变化？
		opponent_winning_move := next_state.Find_Winning_Move(opponent)
		if opponent_winning_move == nil {
			possible_moveArr = append(possible_moveArr, candidate_move)
		}
	}
	return possible_moveArr
}

// 找到两步棋之后可以确保获胜的函数
func (gs *GameState) Find_Two_Step_Win(nextPlayer Player) *Move {
	opponent := nextPlayer.Other()
	for _, candidate_move := range gs.GetAllLegalMoves() {
		next_state := gs.ApplyMove(candidate_move) // 看看走出这一步，棋局会如何变化？

		good_responses := next_state.Eliminate_Losing_Move(opponent) // 对手能不能做出良好的防御？ 如果不能，就选择这一步
		if len(good_responses) <= 0 {
			return &candidate_move
		}
	}
	return nil // 无论己方选择那个动作，对方都能阻止己方获胜
}
