package aigo

import (
	"fmt"
	"math"
)

// 围棋的胜负判断中国规则
// 因为执黑先行，所以结束时黑棋要比白棋多7.5目才算赢。

// 围棋游戏结果判定类
type GameResult struct {
	B    int     // 黑棋气数
	W    int     // 白棋
	KOMI float64 //
}

// 赢家是谁？
func (gr *GameResult) Winner() Player {
	if float64(gr.B) > float64(gr.W)+gr.KOMI {
		return Black
	}
	return White
}

// 赢多少子？
func (gr *GameResult) WinningMargin() float64 {
	w := float64(gr.W) + gr.KOMI
	return math.Abs(float64(gr.B) - w)
}

// 序列化字符串
func (gr *GameResult) String() string {
	w := float64(gr.W) + gr.KOMI
	if float64(gr.B) > w {
		return fmt.Sprintf("黑胜%.1f子", (float64(gr.B) - w))
	}
	return fmt.Sprintf("白胜%.1f子", (w - float64(gr.B)))
}
