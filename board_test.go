package aigo

import (
	"testing"
)

func TestCapture(t *testing.T) {
	var err error
	board := NewBoard(19, 19)
	err = board.PlaceStone(Black, Point{2, 2})
	if err != nil {
		t.Error(err)
	}
	err = board.PlaceStone(White, Point{1, 2})
	if err != nil {
		t.Error(err)
	}

	p := board.Get(Point{2, 2})
	if p == nil {
		t.Errorf("%v应该有棋子1!", Point{2, 2})
	}
	if p != nil && *p != Black {
		t.Errorf("%v应该是黑子1!", Point{2, 2})
	}

	err = board.PlaceStone(White, Point{2, 1})
	if err != nil {
		t.Error(err)
	}
	p = board.Get(Point{2, 2})
	if p == nil {
		t.Errorf("%v应该有棋子2!", Point{2, 2})
	}
	if p != nil && *p != Black {
		t.Errorf("%v应该是黑子2!", Point{2, 2})
	}
	err = board.PlaceStone(White, Point{2, 3})
	if err != nil {
		t.Error(err)
	}
	p = board.Get(Point{2, 2})
	if p == nil {
		t.Errorf("%v应该有棋子!3", Point{2, 2})
	}
	if p != nil && *p != Black {
		t.Errorf("%v应该是黑子3!", Point{2, 2})
	}
	err = board.PlaceStone(White, Point{3, 2})
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(board.String())
	p = board.Get(Point{2, 2})
	if p != nil {
		t.Errorf("%v应该黑子被提了 4!", Point{2, 2})
	}
}

func TestCaptureTwoStones(t *testing.T) {
	var err error
	board := NewBoard(19, 19)

	if board.PlaceStone(Black, Point{2, 2}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(Black, Point{2, 3}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(White, Point{1, 2}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(White, Point{1, 3}) != nil {
		t.Error(err)
	}

	p := board.Get(Point{2, 2})
	if p == nil || *p != Black {
		t.Errorf("位置%v棋子%v不对", Point{2, 2}, p)
	}
	p = board.Get(Point{2, 3})
	if p == nil || *p != Black {
		t.Errorf("位置%v棋子%v不对", Point{2, 3}, p)
	}

	if board.PlaceStone(White, Point{3, 2}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(White, Point{3, 3}) != nil {
		t.Error(err)
	}

	p = board.Get(Point{2, 2})
	if p == nil || *p != Black {
		t.Errorf("位置%v棋子%v不对", Point{2, 2}, p)
	}
	p = board.Get(Point{2, 3})
	if p == nil || *p != Black {
		t.Errorf("位置%v棋子%v不对", Point{2, 3}, p)
	}

	if board.PlaceStone(White, Point{2, 1}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(White, Point{2, 4}) != nil {
		t.Error(err)
	}

	p = board.Get(Point{2, 2})
	if p != nil {
		t.Errorf("位置%v棋子%v不对", Point{2, 2}, p)
	}
	p = board.Get(Point{2, 3})
	if p != nil {
		t.Errorf("位置%v棋子%v不对", Point{2, 3}, p)
	}
}

// 测试不是自杀
func TestCaptureIsNotSuicide(t *testing.T) {
	var err error
	board := NewBoard(19, 19)

	if board.PlaceStone(Black, Point{1, 1}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(Black, Point{2, 2}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(Black, Point{1, 3}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(White, Point{2, 1}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(White, Point{1, 2}) != nil { // 提子 1,1
		t.Error(err)
	}

	p := board.Get(Point{1, 1})
	if p != nil {
		t.Errorf("位置%v棋子%v不对", Point{1, 1}, p)
	}

	p = board.Get(Point{2, 1})
	if p == nil || *p != White {
		t.Errorf("位置%v棋子%v不对", Point{2, 1}, p)
	}

	p = board.Get(Point{1, 2})
	if p == nil || *p != White {
		t.Errorf("位置%v棋子%v不对", Point{1, 2}, p)
	}
}

// 测试减少气
func TestRemoveLiberties(t *testing.T) {
	var err error
	board := NewBoard(5, 5)

	if board.PlaceStone(Black, Point{3, 3}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(White, Point{2, 2}) != nil {
		t.Error(err)
	}

	whiteStoneGroup := board.GetStoneGroup(Point{2, 2})
	if PointSliceEqualBCE(whiteStoneGroup.Liberties, []Point{{2, 3}, {2, 1}, {1, 2}, {3, 2}}) {
		t.Errorf("%v != %v", whiteStoneGroup.Liberties, []Point{{2, 3}, {2, 1}, {1, 2}, {3, 2}})
	}

	if board.PlaceStone(Black, Point{3, 2}) != nil {
		t.Error(err)
	}

	whiteStoneGroup = board.GetStoneGroup(Point{2, 2})
	if PointSliceEqualBCE(whiteStoneGroup.Liberties, []Point{{2, 3}, {2, 1}, {1, 2}}) {
		t.Errorf("%v != %v", whiteStoneGroup.Liberties, []Point{{2, 3}, {2, 1}, {1, 2}})
	}
}

// https://www.jianshu.com/p/80f5f5173fca
func PointSliceEqualBCE(a, b []Point) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// 测试在角算气
func TestEmptyTriangle(t *testing.T) {
	var err error
	board := NewBoard(5, 5)

	if board.PlaceStone(Black, Point{1, 1}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(Black, Point{1, 2}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(Black, Point{2, 2}) != nil {
		t.Error(err)
	}
	if board.PlaceStone(White, Point{2, 1}) != nil {
		t.Error(err)
	}

	blackStoneGroup := board.GetStoneGroup(Point{1, 1})
	if PointSliceEqualBCE(blackStoneGroup.Liberties, []Point{{3, 2}, {2, 3}, {1, 3}}) {
		t.Errorf("%v != %v", blackStoneGroup.Liberties, []Point{{3, 2}, {2, 3}, {1, 3}})
	}
}
