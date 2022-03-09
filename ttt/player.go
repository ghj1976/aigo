package ttt

type Player byte

const (
	None Player = iota // 棋盘上没有棋子
	X
	O
)

func (p Player) String() string {
	switch p {
	case O:
		return "O"
	case X:
		return "X"
	default:
		return "None"
	}
}

func (p Player) Other() Player {
	if p == X {
		return O
	} else {
		return X
	}
}
