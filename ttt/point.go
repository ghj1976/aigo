package ttt

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// row column 采用跟Excel一样的定义， B13 先列后行
// row  排 ， 行 一排 水平方向
// column 柱子  列  垂直方向

type Point struct {
	Row, Col uint16
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

func (p Point) String() string {
	return fmt.Sprintf("%s%d", string(COLS[p.Col-1]), p.Row)
}
