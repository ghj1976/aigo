package main

import (
	"fmt"
	"ghj1976/aigo"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	game := aigo.NewGameOfSize(9, 9)
	bot := map[aigo.Player]aigo.IAgent{
		aigo.Black: aigo.RandomBot{},
		aigo.White: aigo.RandomBot{},
	}

	var err error

	for !game.IsOver() {
		time.Sleep(300 * time.Millisecond) // 避免机器人太快，看不清楚

		// https://kuricat.com/gist/goalng-cli-sfk7e
		fmt.Printf("\x1bc")
		// fmt.Printf("\x1b[2J")

		fmt.Println(game.BoardPosition.PrintBoard())

		bot_move := bot[game.PlayerTurn].SelectMove(game)

		fmt.Println(aigo.PrintMove(game.PlayerTurn, bot_move))

		game, err = game.ApplyMove(bot_move)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("结束!")
}
