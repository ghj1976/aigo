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

func main() {
	rand.Seed(time.Now().Unix())

	reader := bufio.NewReader(os.Stdin)

	game := aigo.NewGameOfSize(9, 9)
	bot := aigo.RandonBot{}

	for !game.IsOver() {
		fmt.Printf("\x1bc") // 清屏
		fmt.Println(game.BoardPosition.PrintBoard())

		var move aigo.Move
		var err error
		if game.PlayerTurn == aigo.Black {
			fmt.Println("清输入:")
			text, ex := reader.ReadString('\n')
			log.Println(text)
			text = strings.ToUpper(text)
			if ex != nil {
				log.Fatalln(ex)
			}
			text = strings.Replace(text, "\n", "", -1)
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
			log.Println(err)
		}
	}
}
