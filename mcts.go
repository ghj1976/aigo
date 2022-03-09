package aigo

import (
	"log"
	"math/rand"
)

type MCTSNode struct {
	game_state      *GameState
	parent          *MCTSNode
	move            *Move
	win_count       map[Player]int
	num_rollouts    int
	children        []*MCTSNode
	unvisited_moves []Move
}

func NewMCTSNode(gs *GameState, parent *MCTSNode, move *Move) *MCTSNode {
	node := &MCTSNode{}
	node.game_state = gs
	node.parent = parent
	node.move = move
	node.win_count = make(map[Player]int)
	node.win_count[Black] = 0
	node.win_count[White] = 0
	node.num_rollouts = 0
	node.children = make([]*MCTSNode, 0)
	node.unvisited_moves = gs.LegalMoves()
	return node
}

// 向树中添加新的子节点
func (node *MCTSNode) add_random_child() *MCTSNode {
	index := rand.Intn(len(node.unvisited_moves))
	new_move := node.unvisited_moves[index]
	new_game_state, err1 := node.game_state.ApplyMove(new_move)
	if err1 != nil {
		log.Fatal(err1)
	}
	new_node := NewMCTSNode(new_game_state, node, &new_move)
	node.children = append(node.children, new_node)
	return new_node
}

// 更新推演统计信息
func (node *MCTSNode) record_win(winner Player) {
	node.win_count[winner]++
	node.num_rollouts++
}

// 检测当前棋局中是否还有合法动作尚未添加到树中
func (node *MCTSNode) can_add_child() bool {
	return len(node.unvisited_moves) > 0
}

// 检测是否达到了终盘。如果已经达到终盘，就不能继续进行搜索了。
func (node *MCTSNode) is_terminal() bool {
	return node.game_state.IsOver()
}

// 返回某一方在推演中获胜的比率。
func (node *MCTSNode) winning_frac(player Player) float64 {
	return float64(node.win_count[player]) / float64(node.num_rollouts)
}
