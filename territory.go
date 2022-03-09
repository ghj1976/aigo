package aigo

// 地盘判断
type Territory struct {
	NumBlackTerritory int     // 黑色地盘
	NumWhiteTerritory int     // 白色地盘
	NumBlackStones    int     // 黑子已占位置数量
	NumWhiteStones    int     // 白子已占位置数量
	NumDame           int     // 中性点数量
	DamePoints        []Point // 中性点
}
