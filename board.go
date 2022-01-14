package aigo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"strings"
)

// 棋盘类，defined by width, height, and a map of stone groups.
type Board struct {
	Width, Height uint16
	stoneMap      map[Point]*StoneGroup // 从棋盘点映射到棋链的map对象
	hash          int64                 // 使用 Zobrist哈希 来增强劫争判断用的
}

// 构造一个指定长宽的棋盘
func NewBoard(w uint16, h uint16) *Board {
	m := make(map[Point]*StoneGroup, w*h)
	return &Board{w, h, m, EmptyBoardHashCode}
}

// Method returns a deep copy of a Board struct.
func (b *Board) Copy() *Board {
	nb := NewBoard(b.Width, b.Height)
	for _, e := range b.GetAllStoneGroups() {
		csg := e.Copy()
		for _, p := range csg.Stones {
			nb.stoneMap[p] = csg
		}
	}
	return nb
}

// Method returns a boolean comparison of two Board structs.
// 比较两个棋盘
func (b *Board) Equal(c *Board) bool {
	if c == nil {
		return false
	}
	if b.Width != c.Width || b.Height != c.Height {
		return false
	}
	bsg, csg := b.GetAllStoneGroups(), c.GetAllStoneGroups() // 所有棋链对比
	if len(bsg) != len(csg) {
		return false
	}
	for _, e1 := range bsg {
		t := false
		for _, e2 := range csg {
			if e1.Equal(e2) { // 比较棋链
				t = true
				break
			}
		}
		if !t {
			return false
		}
	}
	return true
}

// Method implements Stringer interface for Board struct.
func (b *Board) String() string {
	s := fmt.Sprint("Board (", b.Width, "x", b.Height, ") {\n")
	for _, e := range b.GetAllStoneGroups() {
		s += fmt.Sprintln(e)
	}
	return s + "}"
}

// 定义一个字符串变量,其中每个字母代表围棋棋盘的一列。忽略字母I，以免与数字1混淆。
const COLS = "ABCDEFGHJKLMNOPQRST"

func (b *Board) PrintBoard() string {
	bbuf := strings.Builder{}

	for row := int(b.Height); row > 0; row-- {
		bbuf.WriteString(fmt.Sprintf("%02d ", row))
		for col := 1; col <= int(b.Width); col++ {
			pp, ex := b.stoneMap[Point{Col: uint16(col), Row: uint16(row)}]
			if !ex {
				bbuf.WriteString(".") // 空白位置
			} else if pp == nil {
				log.Printf("R%dC%d nil ", row, col)
			} else {
				if pp.Color == Black {
					bbuf.WriteString("X") // 黑棋
				} else {
					bbuf.WriteString("O") // 白棋
				}
			}
		}
		bbuf.WriteString("\r\n")
	}
	bbuf.WriteString("   ")
	for i := 0; i < int(b.Height); i++ {
		bbuf.WriteByte(COLS[i])
	}
	bbuf.WriteString("\r\n")

	return bbuf.String()
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

// 返回棋盘某个交叉点的内容
// 如果已经落子，返回它的棋链对象
// 否则返回 nil
func (b *Board) GetStoneGroup(p Point) *StoneGroup {
	if sg, e := b.stoneMap[p]; e {
		return sg
	}
	return nil
}

// 返回棋盘某个位置的内容
// nil 空
// 否则 是具体的棋子颜色
func (b *Board) Get(p Point) *Player {
	sg := b.GetStoneGroup(p)
	if sg == nil {
		return nil
	}
	return &sg.Color
}

// Method returns a slice of pointers to all unique StoneGroups.
// 返回排重过的棋链列表
func (b *Board) GetAllStoneGroups() []*StoneGroup {
	v := []*StoneGroup{}
	for _, e := range b.stoneMap {
		if e == nil {
			continue
		}
		f := false // 是否已经记录了某个棋链？
		for _, vsg := range v {
			if vsg == e {
				f = true
				break
			}
		}
		if !f {
			v = append(v, e)
		}
	}
	return v
}

// Method returns a new Board with a Player stone played at a Point.
// 指定位置下棋
// 这个函数不限制下棋顺序，可以连续同色下棋，方便让棋、调试等场景
// 如果没气了，棋子会被自动提走
func (b *Board) PlaceStone(turn Player, p Point) error {
	// error checking
	if !b.IsOnGrid(p) { // 是否在棋盘上
		return errors.New("given point is not within the board")
	}
	if b.GetStoneGroup(p) != nil { // 指定位置有棋子了
		return errors.New("given point on the board is already occupied")
	}

	// initialize utilities
	// 棋链列表排重增加
	add_group := func(list []*StoneGroup, item *StoneGroup) []*StoneGroup {
		for _, e := range list {
			if e == item {
				return list
			}
		}
		return append(list, item)
	}
	adjacent_same_color := []*StoneGroup{}
	adjacent_opposite_color := []*StoneGroup{}
	adjacent_liberties := []Point{}

	// collect information about neighbors
	for _, e := range p.Neighbors() {
		if !b.IsOnGrid(e) {
			continue
		}
		nsg := b.GetStoneGroup(e)
		if nsg == nil { // 这个位置是空位， 加气
			adjacent_liberties = append(adjacent_liberties, e)
			continue
		}
		if nsg.Color == turn {
			adjacent_same_color = add_group(adjacent_same_color, nsg)
		} else {
			adjacent_opposite_color = add_group(adjacent_opposite_color, nsg)
		}
	}

	// apply game logic
	newsg := StoneGroup{turn, []Point{p}, adjacent_liberties}

	for _, e := range adjacent_same_color { // 同样颜色的增加棋链
		err := newsg.MergeIn(e)
		if err != nil {
			return err
		}
	}
	for _, e := range newsg.Stones { // 修改棋子到棋链的映射关系
		b.stoneMap[e] = &newsg
	}

	// 使用 Zobrist哈希 后，增加的代码
	// 对棋盘应用这个交叉点与棋子颜色所对应的哈希值
	b.hash ^= BoardPointHashCode[NewBoardPoint(p.Row, p.Col, turn)]

	for _, e := range adjacent_opposite_color { // 不同颜色的减少气
		err := e.RemoveLiberty(p)
		if err != nil {
			return err
		}
	}
	for _, e := range adjacent_opposite_color { // 如果气数为0， 提取棋子
		if e.NumLiberties() == 0 {
			err := b.removeStones(e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Method removes a StoneGroup from the board.
func (b *Board) removeStones(sg *StoneGroup) error {
	for _, e := range sg.Stones {
		if b.stoneMap[e] != sg {
			return errors.New("StoneGroup to be removed does not match board state")
		}
	}
	for _, e := range sg.Stones {
		for _, p := range e.Neighbors() {
			if !b.IsOnGrid(p) {
				continue
			}
			nsg := b.GetStoneGroup(p)
			if nsg != nil && nsg != sg {
				// 加气可以不成功
				nsg.AddLiberty(e)
			}
		}
		b.stoneMap[e] = nil

		// 在 Zobrist哈希中，需要通过逆应用这步动作的哈希值来实现提子
		b.hash ^= BoardPointHashCode[NewBoardPoint(e.Row, e.Col, sg.Color)]

	}
	return nil
}

// 指定的位置，对某个棋子来说是否是眼
// 往自己的眼下棋是不允许的。
// 所有的相邻交叉点以及4个对角相邻点中有3个以上都是己方的棋子才算眼
func (b *Board) IsPointAnEye(p Point, color Player) bool {
	sg := b.GetStoneGroup(p)
	if sg != nil {
		return false // 眼必须是空点
	}

	// 四个相邻的点，都必须是己方的棋
	for _, neighbor := range p.Neighbors() {
		if b.IsOnGrid(neighbor) {
			nsg := b.GetStoneGroup(neighbor)
			if nsg == nil {
				return false // 邻居是空点，不算眼
			}
			if nsg.Color != color {
				return false // 不是自己的棋，不算眼
			}
		}
	}

	// 四个对角线的点，至少有三个是自己的，才算点
	friendly_corners := 0  // 棋盘内的对角线点
	off_board_corners := 0 // 棋盘外的对角线点
	corners := []Point{
		{Row: uint16(p.Row - 1), Col: uint16(p.Col - 1)},
		{Row: p.Row - 1, Col: p.Col + 1},
		{Row: p.Row + 1, Col: p.Col - 1},
		{Row: p.Row + 1, Col: p.Col + 1},
	}
	for _, corner := range corners {
		if b.IsOnGrid(corner) {
			csg := b.GetStoneGroup(corner)
			if csg != nil {
				if csg.Color == color {
					friendly_corners++
				}
			}
		} else {
			off_board_corners++
		}
	}
	if off_board_corners > 0 {
		// 空点在边缘或角落
		return off_board_corners+friendly_corners == 4
	} else {
		// 空点在棋盘里
		return friendly_corners >= 3
	}
}

// 检查棋子是否在棋盘上
// 棋盘坐标是从1开始的，包含1
func (b *Board) IsOnGrid(p Point) bool {
	b1 := (1 <= p.Row) && (p.Row <= b.Height)
	b2 := (1 <= p.Col) && (p.Col <= b.Width)
	return b1 && b2
}

// 返回 Zobrist 哈希值
func (b *Board) GetZobristHash() int64 {
	return b.hash
}
