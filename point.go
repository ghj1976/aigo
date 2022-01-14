package aigo

import (
	"log"
	"strconv"
	"strings"
)

// A point on the board, defined by row and column.
// uint16 范围 0 ~ 65535 (6万多)
type Point struct {
	Row, Col uint16
}

// 周围上下左右四个位置的坐标
func (p Point) Neighbors() [4]Point {
	r := [4]Point{}
	r[0] = Point{p.Row, p.Col - 1}
	r[1] = Point{p.Row + 1, p.Col}
	r[2] = Point{p.Row, p.Col + 1}
	r[3] = Point{p.Row - 1, p.Col}
	return r
}

// 把手工输入的 B17 这样的转换成 2,17
// 人机对弈场景使用
func PointFromCoords(coords string) *Point {
	col := uint16(strings.Index(COLS, string(coords[0])) + 1)
	row, e := strconv.Atoi(coords[1:])
	if e != nil {
		log.Panicln(e)
		return nil
	}
	return &Point{Row: uint16(row), Col: col}

}
