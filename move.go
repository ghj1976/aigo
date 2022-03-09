package aigo

import (
	"fmt"
	"strings"
)

// A move: either play (at given point), pass, or resign.
// 美国围棋协会（American GO Association，AGA）的惯例，我们使用术语动作（move）来表示这3种行动中的任何一个，而用落子（play）表示落下一颗棋子。
type Move struct {
	Pnt      Point // 下棋的位置
	IsPlay   bool  // 正常下棋
	IsPass   bool  // 跳过
	IsResign bool  // 认输
}

// Constructs a play Move with passed point.
// 落子
func NewPlay(p Point) Move {
	return Move{p, true, false, false}
}

// Constructs a pass Move.
// 跳过
func NewPass() Move {
	return Move{Point{}, false, true, false}
}

// Constructs a resign Move.
// 认输
func NewResign() Move {
	return Move{Point{}, false, false, true}
}

// Method implements Stringer interface for Move struct.
func (m Move) String() string {
	if m.IsPlay && !(m.IsPass || m.IsResign) {
		return fmt.Sprint("Play{", m.Pnt.Row, m.Pnt.Col, "}")
	} else if m.IsPass && !m.IsResign {
		return "Pass 跳过"
	} else if m.IsResign {
		return "Resign 认输"
	} else {
		return "Invalid"
	}
}

func (m Move) StringChessRecord() string {
	if m.IsPlay && !(m.IsPass || m.IsResign) {
		return fmt.Sprintf("%s%d", string(COLS[m.Pnt.Col-1]), m.Pnt.Row)
	} else {
		return m.String()
	}
}

// 打印下棋动作
func PrintMove(p Player, m Move) string {
	bbuf := strings.Builder{}
	bbuf.WriteString(p.String())
	bbuf.WriteString(" ")
	bbuf.WriteString(m.StringChessRecord())
	return bbuf.String()
}
