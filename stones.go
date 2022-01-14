package aigo

import (
	"errors"
	"fmt"
)

// 棋链 一组同色相连的棋子称为一条棋链
type StoneGroup struct {
	Color     Player  // 玩家颜色
	Stones    []Point // 棋链包含的棋子位置 stone 石头
	Liberties []Point // 棋链有气的位置  liberty 自由
}

// 返回深度copy的对象（deep copy）
func (sg *StoneGroup) Copy() *StoneGroup {
	ns := make([]Point, sg.NumStones())
	copy(ns, sg.Stones)
	nl := make([]Point, sg.NumLiberties())
	copy(nl, sg.Liberties)
	return &StoneGroup{sg.Color, ns, nl}
}

// 比较两个棋链
func (sg *StoneGroup) Equal(c *StoneGroup) bool {
	if c == nil {
		return false
	}
	b1, b2, b3 := c.Color != sg.Color, c.NumStones() != sg.NumStones(), c.NumLiberties() != sg.NumLiberties()
	if b1 || b2 || b3 { // 颜色，棋子数、气数不一样，则不是一个棋链
		return false
	}
	for _, e := range sg.Stones { // 棋子位置比较
		if _, t := contains(c.Stones, e); !t {
			return false
		}
	}
	for _, e := range sg.Liberties { // 气位置比较
		if _, t := contains(c.Liberties, e); !t {
			return false
		}
	}
	return true
}

// Method implements Stringer interface for StoneGroup struct.
func (sg *StoneGroup) String() string {
	s1, s2, s3 := "StoneGroup{\n -Color: ", " -Stones: ", "\n -Liberties: "
	c := "Black\n"
	if sg.Color != Black {
		c = "White\n"
	}
	return fmt.Sprint(s1, c, s2, sg.Stones, s3, sg.Liberties, " }")
}

// 棋链的气数
func (sg *StoneGroup) NumLiberties() int {
	return len(sg.Liberties)
}

// 棋链的棋子数
func (sg *StoneGroup) NumStones() int {
	return len(sg.Stones)
}

// 给棋链增加气
func (sg *StoneGroup) AddLiberty(p Point) error {
	if _, b := contains(sg.Liberties, p); b {
		return errors.New("stoneGroup already contains the given liberty")
	} else {
		sg.Liberties = append(sg.Liberties, p)
		return nil
	}
}

// 给棋链删除气
func (sg *StoneGroup) RemoveLiberty(p Point) error {
	if i, b := contains(sg.Liberties, p); !b { // i 棋链所在数组的位置标号
		return errors.New("stoneGroup doesn't have the given liberty")
	} else {
		s, l := sg.Liberties, sg.NumLiberties()
		s[l-1], s[i] = s[i], s[l-1] // 跟最后一个互换
		sg.Liberties = s[:l-1]      // 删除换后的最后一个
		return nil
	}
}

// 合并棋链，包含sg和mg这两个的新棋链。
// 多个棋链时，需要依次合并，这里只考虑2条棋链之间的关系，不考虑三条。
func (sg *StoneGroup) MergeIn(mg *StoneGroup) error {
	if sg.Color != mg.Color {
		return errors.New("cannot merge StoneGroup of different player color")
	}
	cs := sg.Stones
	for _, e := range mg.Stones { // 合并棋子
		if _, b := contains(cs, e); !b { // 如果没有增加
			cs = append(cs, e)
		}
	}
	cl := sg.Liberties
	for _, e := range mg.Liberties {
		if _, b1 := contains(cl, e); !b1 { // 不是已经存在的气
			if _, b2 := contains(cs, e); !b2 { // 原先的气没有被放棋子
				cl = append(cl, e)
			}
		}
	}
	sg.Stones = cs
	sg.Liberties = cl
	return nil
}

// 公共函数 检查点列表中是否包含某个点
func contains(s []Point, e Point) (int, bool) {
	for i, a := range s {
		if a == e {
			return i, true
		}
	}
	return -1, false
}
