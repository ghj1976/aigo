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
				log.Printf("Row:%d,Col:%d is nil", row, col)
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
// 否则 是具体的棋子颜色
func (b *Board) Get(p Point) Player {
	sg := b.GetStoneGroup(p)
	if sg == nil {
		return None
	}
	return sg.Color
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
		delete(b.stoneMap, e)
		// b.stoneMap[e] = nil 绝对不能这么做，会有问题的。

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

// 评估各方领土
// 假设棋盘上所有死棋都被提走了，然后开始计算胜负
func (b *Board) EvaluateTerritory() *Territory {
	visited := make(map[Point]bool)
	territory_map := make(map[Point]string)
	for r := uint16(1); r <= b.Height; r++ {
		for c := uint16(1); c <= b.Width; c++ {
			p := Point{Row: r, Col: c}

			if _, ex1 := territory_map[p]; ex1 { // 已经分析过这个点了
				continue
			}

			stone := b.Get(p)
			if stone != None {
				territory_map[p] = stone.String()
			} else {
				group, neighbors := b.scoringCollectRegion(p, visited)
				fill_with := ""
				if len(neighbors) == 1 {
					fill_with = fmt.Sprintf("territory_%s", neighbors[0].String())
				} else {
					fill_with = "dame"
				}
				for _, pos := range group {
					territory_map[pos] = fill_with
				}
			}
		}
	}

	territory := &Territory{}

	for point, status := range territory_map {

		switch status {
		case "Black":
			territory.NumBlackStones++
		case "White":
			territory.NumWhiteStones++
		case "territory_Black":
			territory.NumBlackTerritory++
		case "territory_White":
			territory.NumWhiteTerritory++
		case "dame":
			territory.NumDame++
			territory.DamePoints = append(territory.DamePoints, point)
		default:
			log.Fatalf("错误的状态 %v", status)
		}
	}
	return territory
}

// 递归收集棋盘空点位置归属
// start_pos 开始位置
// visited 用于判断该位置是否已经收集
// Find the contiguous section of a board containing a point. Also identify all the boundary points.
// 只有完全被一个棋链包围的才算自己的气
// all_points 周围所有跟起始点都是空的点集合，这些点都会被递归调用到
// all_borders 这些空白点的边界情况， （非起始点同色的）
func (b *Board) scoringCollectRegion(start_pos Point, visited map[Point]bool) (all_points []Point, all_borders []Player) {

	_, ex1 := visited[start_pos]
	if ex1 { // 已经分析过这个位置了
		return []Point{}, []Player{}
	}

	all_points = []Point{start_pos}
	all_borders = []Player{}
	visited[start_pos] = true
	here := b.Get(start_pos)
	deltas := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, delta := range deltas {
		delta_r := delta[0]
		delta_c := delta[1]

		next_p := Point{Row: start_pos.Row + uint16(delta_r), Col: start_pos.Col + uint16(delta_c)}
		if !b.IsOnGrid(next_p) { // 超出棋盘范围
			continue
		}

		neighbor := b.Get(next_p)

		if neighbor == here { // 这个位置是空位，需要判断归属谁？

			// 只有跟起始点同色的才会递归一直找下去。
			points, borders := b.scoringCollectRegion(next_p, visited)
			all_points = pointArrMerge(all_points, points)
			all_borders = playerArrMerge(all_borders, borders)
		} else {
			// 增加空白区域的周边信息
			all_borders = playerArrUpdate(all_borders, neighbor)
		}
	}

	return all_points, all_borders

}

// 如果不存在，则增加，存在则不变
// TODO go 2.0 时优化下
func playerArrUpdate(arr []Player, p Player) []Player {
	find := false
	for _, n := range arr {
		if n == p {
			find = true
		}
	}
	if find {
		return arr
	}
	return append(arr, p)
}

// 两个数组的合并， 重复的作为一项
func playerArrMerge(arr1, arr2 []Player) []Player {
	for _, n := range arr2 {
		arr1 = playerArrUpdate(arr1, n)
	}
	return arr1
}

func pointArrUpdate(arr []Point, p Point) []Point {
	find := false
	for _, n := range arr {
		if n == p {
			find = true
		}
	}
	if find {
		return arr
	}
	return append(arr, p)
}

func pointArrMerge(arr1, arr2 []Point) []Point {
	for _, n := range arr2 {
		arr1 = pointArrUpdate(arr1, n)
	}
	return arr1
}
