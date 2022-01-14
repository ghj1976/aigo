package main

import (
	"fmt"
	"ghj1976/aigo"
	"log"
)

// 3.3 的目的是验证各项功能都具备，是要写单元测试的。
func main() {
	var err error
	game := aigo.NewGameOfSize(9, 9)
	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 3, Col: 4}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 2, Col: 4}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 3, Col: 3}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 2, Col: 3}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 3, Col: 2}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 2, Col: 2}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 3, Col: 1}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 2, Col: 1}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 2, Col: 5}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 1, Col: 2}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 1, Col: 5}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 1, Col: 4}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(aigo.NewPlay(aigo.Point{Row: 5, Col: 5}))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(game)
	fmt.Println(game.BoardPosition.PrintBoard())

	if game.BoardPosition.IsPointAnEye(aigo.Point{Row: 1, Col: 1}, aigo.White) {
		fmt.Printf("%v下在眼（%v）上了", aigo.White, aigo.Point{Row: 1, Col: 1})
	}

	if game.BoardPosition.IsPointAnEye(aigo.Point{Row: 1, Col: 1}, aigo.Black) {
		fmt.Printf("%v下在眼（%v）上了", aigo.Black, aigo.Point{Row: 1, Col: 1})
	}

	// game, err = game.ApplyMove(aigo.White, aigo.NewPlay(aigo.Point{Row: 1, Col: 1}))
	// if err != nil {
	// 	log.Fatalln(err)
	// }

}
