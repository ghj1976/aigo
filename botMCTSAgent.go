package aigo

import (
	"log"
	"math"
)

//
type MCTSAgent struct {
	num_rounds  int //
	temperature float64
}

func NewMCTSAgent(numrounds int, temperature float64) *MCTSAgent {
	bot := &MCTSAgent{num_rounds: numrounds, temperature: temperature}
	return bot
}

func (bot *MCTSAgent) SelectMove(gs *GameState) Move {
	root := NewMCTSNode(gs, nil, nil)

	for i := 0; i < bot.num_rounds; i++ {
		node := root
		for !node.can_add_child() && !node.is_terminal() {
			node = bot.select_child(*node)
		}

		// Add a new child node into the tree.
		if node.can_add_child() {
			node = node.add_random_child()
		}

		log.Printf("curr %d :%v ", i, node.move)
		winner := bot.simulate_random_game(node.game_state)

		for node != nil {
			node.record_win(winner)
			node = node.parent
		}
	}

	var best_move *Move
	best_pct := -1.0
	for _, child := range root.children {
		child_pct := child.winning_frac(gs.PlayerTurn)
		if child_pct > best_pct {
			best_pct = child_pct
			best_move = child.move
		}
	}
	log.Println(len(root.children))
	return *best_move
}

// 使用 搜索树置信区间上界公式 找一个应该探索的节点
func (bot *MCTSAgent) select_child(node MCTSNode) *MCTSNode {
	// 搜索树置信区间上界公式（upper confidence bound for trees formula，简称为UCT公式）
	total_rollouts := 0
	for _, child := range node.children {
		total_rollouts += child.num_rollouts
	}
	log_rollouts := math.Log(float64(total_rollouts))

	best_score := -1.0
	var best_child *MCTSNode

	for _, child := range node.children {
		win_percentage := child.winning_frac(node.game_state.PlayerTurn)
		exploration_factor := math.Sqrt(log_rollouts / float64(child.num_rollouts))
		uct_score := win_percentage + bot.temperature*exploration_factor
		if uct_score > best_score {
			best_score = uct_score
			best_child = child
		}
	}
	return best_child
}

// 随机模拟一盘游戏
func (bot *MCTSAgent) simulate_random_game(gs *GameState) Player {
	bots := map[Player]IAgent{
		White: NewFastRandomBot(),
		Black: NewFastRandomBot(),
	}
	var err error
	for !gs.IsOver() {
		bot_move := bots[gs.PlayerTurn].SelectMove(gs)
		gs, err = gs.ApplyMove(bot_move)
		if err != nil {
			log.Fatal(err)
		}

	}
	return gs.Winner()
}
