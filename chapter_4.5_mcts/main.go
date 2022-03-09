package main

import (
	"bufio"
	"fmt"
	"ghj1976/aigo"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	BOARD_SIZE = 5 // 	棋盘大小
)

func main() {
	rand.Seed(time.Now().Unix())
	reader := bufio.NewReader(os.Stdin)

	game := aigo.NewGameOfSize(BOARD_SIZE, BOARD_SIZE)

	bot := aigo.NewMCTSAgent(500, 1.4)

	for !game.IsOver() {
		// fmt.Printf("\x1bc") // 清屏
		fmt.Println(game.BoardPosition.PrintBoard())

		var move aigo.Move
		var err error
		if game.PlayerTurn == aigo.Black {

			text := ""
			for {
				fmt.Print("请输入:")
				var ex error
				text, ex = reader.ReadString('\n')
				log.Println(text)
				text = strings.ToUpper(text)
				if ex != nil {
					log.Fatalln(ex)
				}
				text = strings.Replace(text, "\n", "", -1)
				if len(text) == 2 {
					break
				}
			}

			p := aigo.PointFromCoords(text)
			if p != nil {
				move = aigo.NewPlay(*p)
			}
		} else {
			move = bot.SelectMove(game)
		}

		fmt.Println(aigo.PrintMove(game.PlayerTurn, move))

		game, err = game.ApplyMove(move)
		if err != nil {
			log.Printf("game.ApplyMove %v %v\r\n", move, err)
		}
	}
}
