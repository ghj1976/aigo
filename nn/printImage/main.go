package main

import (
	"ghj1976/aigo/nn/mnist"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {

	dataSet, err := mnist.ReadTrainSet("../mnist")
	if err != nil {
		log.Fatal(err)
	}

	imCols := 28
	imRows := 28

	rect := image.Rect(0, 0, imCols, imRows)

	rgba := image.NewNRGBA(rect)

	log.Println(dataSet.Data[0].Digit)
	for dy := 0; dy < imCols; dy++ {
		for dx := 0; dx < imRows; dx++ {
			rgba.Set(dy, dx, color.Gray{dataSet.Data[0].Image[dx][dy]})
		}
	}

	fIm, err := os.Create("a0.png")

	if nil != err {
		log.Fatal(err)
	}

	err = png.Encode(fIm, rgba)

	if nil != err {
		log.Fatal(err)
	}

}
