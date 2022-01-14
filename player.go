package aigo

type Player byte

// 在围棋对弈中，黑方与白方轮流落子，因此可以用enum类型来表示棋子的颜色。
const (
	None Player = iota // 棋盘上没有棋子
	Black
	White
)

func (p Player) String() string {
	if p == Black {
		return "Black"
	} else if p == None {
		return "None"
	} else {
		return "White"
	}
}

// 每一回合棋手落子之后，调用Player实例上的Other方法来进行切换棋手。
func (p Player) Other() Player {
	if p == Black {
		return White
	} else { // 如果是空玩家，那就认为是刚开始，黑棋先行
		return Black
	}
}

func (p Player) GetNone() Player {
	return None
}

func (p Player) GetBlack() Player {
	return Black
}

func (p Player) GetWhite() Player {
	return White
}
