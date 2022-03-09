package ttt

import (
	"fmt"
	"strings"
)

const (
	BOARD_SIZE = 3 // 棋盘大小
)

type Board struct {
	grid map[Point]Player
}

func NewBoard() *Board {
	b := &Board{}
	b.grid = make(map[Point]Player)
	return b
}

// 检查棋子是否在棋盘上
// 棋盘坐标是从1开始的，包含1
func (b *Board) IsOnGrid(p Point) bool {
	b1 := (1 <= p.Row) && (p.Row <= BOARD_SIZE)
	b2 := (1 <= p.Col) && (p.Col <= BOARD_SIZE)
	return b1 && b2
}

func (b *Board) Get(p Point) Player {
	pp, ok := b.grid[p]
	if !ok {
		return None
	}
	return pp
}

func (b *Board) Place(p Player, point Point) {
	if !b.IsOnGrid(point) {
		panic(fmt.Sprintf("%v不在棋盘上", point))
	}
	if b.Get(point) != None {
		panic(fmt.Sprintf("%v 位置已经有棋%v了", point, b.Get(point)))
	}

	b.grid[point] = p
}

// 是否满足3
func (b *Board) Has_3_in_a_row(p Player) bool {
	if p == None {
		return false
	}

	// 每列 col 是否满3
	for col := 1; col <= BOARD_SIZE; col++ {
		if b.grid[Point{1, uint16(col)}] == None {
			continue
		}
		if b.grid[Point{1, uint16(col)}] == p &&
			(b.grid[Point{2, uint16(col)}] == b.grid[Point{3, uint16(col)}]) &&
			(b.grid[Point{1, uint16(col)}] == b.grid[Point{2, uint16(col)}]) {
			return true
		}
	}

	// 每行 row 是否满3
	for row := 1; row <= BOARD_SIZE; row++ {
		if b.grid[Point{uint16(row), 1}] == None {
			continue
		}
		if b.grid[Point{uint16(row), 1}] == p &&
			(b.grid[Point{uint16(row), 2}] == b.grid[Point{uint16(row), 3}]) &&
			(b.grid[Point{uint16(row), 1}] == b.grid[Point{uint16(row), 2}]) {
			return true
		}
	}

	// 交叉线 \
	if b.grid[Point{3, 1}] == p &&
		(b.grid[Point{2, 2}] == b.grid[Point{1, 3}]) &&
		(b.grid[Point{3, 1}] == b.grid[Point{2, 2}]) {
		return true
	}

	// 交叉线 /
	if b.grid[Point{1, 1}] == p &&
		(b.grid[Point{2, 2}] == b.grid[Point{1, 1}]) &&
		(b.grid[Point{3, 3}] == b.grid[Point{2, 2}]) {
		return true
	}

	return false
}

// Method returns a deep copy of a Board struct.
func (b *Board) DeepCopy() *Board {
	nb := NewBoard()
	for k, v := range b.grid {
		nb.grid[k] = v
	}
	return nb
}

// 定义一个字符串变量,其中每个字母代表围棋棋盘的一列。忽略字母I，以免与数字1混淆。
const COLS = "ABCDEFGHJKLMNOPQRST"

func (b *Board) PrintBoard() string {
	bbuf := strings.Builder{}

	for row := BOARD_SIZE; row > 0; row-- {
		bbuf.WriteString(fmt.Sprintf("%02d ", row))
		for col := 1; col <= BOARD_SIZE; col++ {
			pp := b.Get(Point{Row: uint16(row), Col: uint16(col)})
			switch pp {
			case X:
				bbuf.WriteString("X")
			case O:
				bbuf.WriteString("O")
			default:
				bbuf.WriteString(".")
			}

		}
		bbuf.WriteString("\r\n")
	}
	bbuf.WriteString("   ")
	for i := 0; i < BOARD_SIZE; i++ {
		bbuf.WriteByte(COLS[i])
	}
	bbuf.WriteString("\r\n")

	return bbuf.String()
}
