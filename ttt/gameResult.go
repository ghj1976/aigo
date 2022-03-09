package ttt

type GameResult byte

const (
	Loss GameResult = iota // 失败
	Draw                   // 平手
	Win                    // 赢
)

// 反转游戏结果
func (gr GameResult) reverse_game_result() GameResult {
	switch gr {
	case Loss:
		return Win
	case Win:
		return Loss
	default:
		return Draw
	}
}
