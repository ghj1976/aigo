package aigo

// 棋盘上点的信息
type BoardPoint struct {
	Point  // 棋盘的位置
	Player // 棋盘上可能没棋子， nil， 也可能有棋子
}

func NewBoardPoint(row, col uint16, player Player) BoardPoint {
	bp := BoardPoint{}
	bp.Row = row
	bp.Col = col
	bp.Player = player
	return bp
}

var (
	BoardPointHashCode2 = MapInitTest()
	EmptyBoardHashCode2 = int64(9181944435492932548)
)

func MapInitTest() map[BoardPoint]int64 {
	mappp := make(map[BoardPoint]int64)
	mappp[NewBoardPoint(1, 1, None)] = 6402364705153495313
	mappp[NewBoardPoint(1, 1, White)] = 444191475187629924
	mappp[NewBoardPoint(1, 1, Black)] = 3180807544946524599

	return mappp
}
