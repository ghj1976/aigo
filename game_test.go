package aigo

import (
	"fmt"
	"log"
	"testing"
)

// 测试连续下棋
func TestSamePlay(t *testing.T) {
	var err error
	game := NewGameOfSize(9, 9)
	game, err = game.ApplyMove(NewPlay(Point{Row: 2, Col: 2}))
	if err != nil {
		log.Fatalln(err)
	}
	game, err = game.ApplyMove(NewPlay(Point{Row: 2, Col: 3}))
	if err != nil {
		log.Fatalln(err)
	}
	game, err = game.ApplyMove(NewPlay(Point{Row: 3, Col: 3}))
	if err != nil {
		log.Fatalln(err)
	}
	game, err = game.ApplyMove(NewPlay(Point{Row: 3, Col: 4}))
	if err != nil {
		log.Fatalln(err)
	}
	game, err = game.ApplyMove(NewPlay(Point{Row: 3, Col: 3}))
	if err != nil {
		log.Println(game.PlayerTurn)
		if err.Error() != "given point on the board is already occupied" {
			t.Error(err)
			t.Errorf("这个位置已经有棋子了!")
		}
	}
	log.Println(game)

}

// 测试气数
func TestLiberties(t *testing.T) {
	var err error
	game := NewGameOfSize(9, 9)
	game, err = game.ApplyMove(NewPlay(Point{Row: 1, Col: 1}))
	if err != nil {
		log.Fatalln(err)
	}

	if len(game.BoardPosition.stoneMap[Point{Row: 1, Col: 1}].Liberties) != 2 {
		log.Println(game)
		fmt.Println(game.BoardPosition.PrintBoard())
		t.Errorf("气数不正确001!")
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 3, Col: 1}))
	if err != nil {
		log.Fatalln(err)
	}
	if len(game.BoardPosition.stoneMap[Point{Row: 3, Col: 1}].Liberties) != 3 {
		log.Println(game)
		fmt.Println(game.BoardPosition.PrintBoard())
		t.Errorf("气数不正确002!")
	}

}

func TestIsPointAnEye(t *testing.T) {
	var err error
	game := NewGameOfSize(9, 9)
	game, err = game.ApplyMove(NewPlay(Point{Row: 3, Col: 4}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 2, Col: 4}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 3, Col: 3}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 2, Col: 3}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 3, Col: 2}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 2, Col: 2}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 3, Col: 1}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 2, Col: 1}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 2, Col: 5}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 1, Col: 2}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 1, Col: 5}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 1, Col: 4}))
	if err != nil {
		log.Fatalln(err)
	}

	game, err = game.ApplyMove(NewPlay(Point{Row: 5, Col: 5}))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(game)
	fmt.Println(game.BoardPosition.PrintBoard())

	if !game.BoardPosition.IsPointAnEye(Point{Row: 1, Col: 1}, White) {
		t.Errorf("%v下在眼（%v）上了", White, Point{Row: 1, Col: 1})
	}

}
