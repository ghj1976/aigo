package main

import (
	"bufio"
	"fmt"
	"ghj1976/aigo/ttt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	human_player := ttt.X

	reader := bufio.NewReader(os.Stdin)
	game := ttt.NewGameState()
	bot := ttt.MinimaxAgent{}

	for !game.IsOver() {
		// fmt.Printf("\x1bc") // 清屏
		fmt.Println(game.BoardPosition.PrintBoard())

		var move ttt.Move
		if game.PlayerTurn == human_player {
			fmt.Print("请输入:")
			text, ex := reader.ReadString('\n')
			log.Println(text)
			text = strings.ToUpper(text)
			if ex != nil {
				log.Fatalln(ex)
			}
			text = strings.Replace(text, "\n", "", -1)
			p := ttt.PointFromCoords(text)
			if p != nil {
				move = ttt.Move{P: *p}
			}

		} else {
			move = bot.SelectMove(*game)
		}
		game = game.ApplyMove(move)
	}

	fmt.Printf("\x1bc") // 清屏
	fmt.Println(game.BoardPosition.PrintBoard())
	winner := game.Winner()
	if winner == ttt.None {
		fmt.Println("平局")
	} else {
		fmt.Printf("Winner:%v\r\n", winner)
	}
}
