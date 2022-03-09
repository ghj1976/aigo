package aigo

import (
	"fmt"
	"log"
)

// A game state, defined by board position, next turn, previous game state, and last move made.
type GameState struct {
	BoardPosition               *Board     // 棋子的布局
	PlayerTurn                  Player     // 当前是那个玩家的回合
	PreviousState               *GameState // 上一回合的游戏状态  在 Zobrist 提速后，这个仍然保留，作为上一步信息的记录
	PreviousZobristHashStateArr []int64    // 提速用的，之前回合的Zobrist哈希数组
	LastMove                    *Move      // 上一步动作
}

// 围棋默认19*19棋盘
func NewGame() *GameState {
	return NewGameOfSize(19, 19)
}

// Constructor function builds a new GameState with Board of passed width and height.
func NewGameOfSize(h uint16, w uint16) *GameState {
	gs := &GameState{}
	gs.BoardPosition = NewBoard(h, w)
	gs.PlayerTurn = Black // 默认是黑棋先行
	gs.PreviousState = nil
	gs.LastMove = nil
	gs.PreviousZobristHashStateArr = []int64{gs.BoardPosition.GetZobristHash()}
	return gs
}

func NewGameState(board *Board, next_player Player, previous *GameState, last_move *Move) *GameState {
	gs := &GameState{}
	gs.BoardPosition = board
	gs.PlayerTurn = next_player
	gs.PreviousState = previous
	gs.LastMove = last_move
	arr, _ := updateInt64Arr(previous.PreviousZobristHashStateArr, board.GetZobristHash())
	gs.PreviousZobristHashStateArr = arr

	return gs
}

// Method implements Stringer interface for GameState struct.
func (gs *GameState) String() string {
	s := fmt.Sprintln("Next turn: ", gs.PlayerTurn, "\nLast move: ", gs.LastMove)
	return fmt.Sprint(s, gs.BoardPosition)
}

// Method returns a new GameState created from applying the given move.
// 执行落子动作后，返回新的 GameState 对象指针
// 不做规则判断
// 下棋的顺序是固定的， 不用传递谁下的这个参数
func (gs *GameState) ApplyMove(m Move) (*GameState, error) {
	var next_board *Board
	if m.IsPlay {
		next_board = gs.BoardPosition.Copy()
		err := next_board.PlaceStone(gs.PlayerTurn, m.Pnt)
		if err != nil {
			return gs, err
		}
	} else {
		next_board = gs.BoardPosition // 跳过或认输场景，棋盘不变
	}
	return NewGameState(next_board, gs.PlayerTurn.Other(), gs, &m), nil
	// return &GameState{next_board, gs.PlayerTurn.Other(), gs, &m}, nil
}

// 判断当前游戏状态是否违反了劫争规则
func (gs *GameState) DoesMoveViolateKo(p Player, move Move) bool {
	if !move.IsPlay {
		return false
	}
	next_board := gs.BoardPosition.Copy()
	next_board.PlaceStone(p, move.Pnt)

	_, ex := containsInt64(gs.PreviousZobristHashStateArr, next_board.hash)
	return ex
	/*
		// 原先没有用 Zobrist的代码
			past_state := gs.PreviousState
			for past_state != nil { // 遍历之前的，
				if past_state.PlayerTurn == p {
					if past_state.BoardPosition == next_board {
						return true
					}
				}

				past_state = past_state.PreviousState
			}
			return false
	*/
}

// 是不是自杀的判断
func (gs *GameState) IsMoveSelfCapture(p Player, move Move) bool {
	if !move.IsPlay { // 跳过和结束不判断
		return false
	}
	nest_board := gs.BoardPosition.Copy()
	nest_board.PlaceStone(p, move.Pnt)
	newSG := nest_board.GetStoneGroup(move.Pnt)
	if newSG == nil { // 棋子已经被提走了，肯定是自杀
		return true
	}
	return newSG.NumLiberties() == 0 // 没被提走，气数为0， 这个分支正常走不到，只是做最后一步检查。
}

// 在给定游戏目前状态下，判断这个动作是否合法？
func (gs *GameState) IsValidMove(move Move) bool {
	if gs.IsOver() {
		return false
	}

	if move.IsPass || move.IsResign { // 认输或者跳过
		log.Println("认输或跳过")
		return true
	}

	if gs.BoardPosition.Get(move.Pnt) != None { // 这个位置已经有棋子了
		// log.Println("已经有棋子")
		return false
	}

	if gs.IsMoveSelfCapture(gs.PlayerTurn, move) { // 出现了自吃， 填自己气的情况
		// log.Println("出现了自吃")
		return false
	}

	if gs.DoesMoveViolateKo(gs.PlayerTurn, move) { // 出现了劫争
		// log.Println("出现了劫争")
		return false
	}

	return true
}

// Method determines if a GameState represents a finished game.
// 是否是决定围棋比赛结束的时机
func (gs *GameState) IsOver() bool {
	lm := gs.LastMove // 之前还有棋手在走，不是结束
	if lm == nil {    // 刚开始会是这样的
		return false
	}
	if lm.IsResign { // 一方主动认输，是结束了
		return true
	}
	slm := gs.PreviousState.LastMove
	if slm == nil { // 刚开始会是这样的
		// panic("gs.PreviousState.LastMove is nil")
		return false
	}
	return lm.IsPass && slm.IsPass // 双方棋手都说跳过不走了，才算结束。
}

// 公共函数 检查点列表中是否包含某个点
func containsInt64(s []int64, e int64) (int, bool) {
	for i, a := range s {
		if a == e {
			return i, true
		}
	}
	return -1, false
}

// 数组增加， 如果存在，不增加，
// 如果不存在，增加
func updateInt64Arr(s []int64, e int64) (arr []int64, ex bool) {
	_, xx := containsInt64(s, e)
	if xx {
		return s, false // 已经有了
	} else {
		arr := append(s, e)
		return arr, true // 原先没有，增加成功
	}

}

// 游戏的赢家
func (gs *GameState) Winner() Player {
	if !gs.IsOver() {
		return None
	}

	if gs.LastMove.IsResign {
		return gs.PlayerTurn
	}

	game_result := gs.ComputeGameResult()
	return game_result.Winner()
}

func (gs *GameState) ComputeGameResult() *GameResult {
	territory := gs.BoardPosition.EvaluateTerritory()
	return &GameResult{
		B:    territory.NumBlackStones + territory.NumBlackTerritory,
		W:    territory.NumWhiteStones + territory.NumWhiteTerritory,
		KOMI: 7.5,
	}
}

// 获得所有合法的可下棋点
func (gs *GameState) LegalMoves() []Move {
	moves := []Move{}
	var r, c uint16
	for r = 1; r <= gs.BoardPosition.Width; r++ {
		for c = 1; c <= gs.BoardPosition.Height; c++ {

			move := NewPlay(Point{Row: r, Col: c})
			if gs.IsValidMove(move) {
				moves = append(moves, move)
			}
		}
	}
	// These two moves are always legal.
	moves = append(moves, NewPass())
	moves = append(moves, NewResign())

	return moves
}

// 棋局评估函数
// 计算棋盘上留存棋子的数量差， 即 将吃掉的棋子相加，然后减去对方吃掉的棋子数量。
// 下一回合轮到谁，就是对谁的棋盘评估结果
// 整个第四章，这个函数就没变过，所以提取到这里
func CaptureDiff(gs *GameState) int {
	black_stones := 0
	white_stones := 0
	var r, c uint16
	for r = 1; r <= gs.BoardPosition.Width; r++ {
		for c = 1; c <= gs.BoardPosition.Height; c++ {
			p := gs.BoardPosition.Get(Point{Row: r, Col: c})
			if p == White {
				white_stones++
			} else if p == Black {
				black_stones++
			}
		}
	}
	diff := black_stones - white_stones
	if gs.PlayerTurn == Black {
		return diff
	}
	return -1 * diff
}
