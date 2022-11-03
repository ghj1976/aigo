package main

import (
	"ghj1976/aigo/nn/mnist"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {

	dataSet, err := mnist.ReadTrainSet("../mnist")
	if err != nil {
		log.Fatal(err)
	}

	digits := [28][28]float64{}
	num := 0
	// 计算数字 8 的平均值
	for _, d := range dataSet.Data {
		if d.Digit != 8 {
			continue
		}

		for x := 0; x < 28; x++ {
			for y := 0; y < 28; y++ {
				digits[x][y] += float64(d.Image[x][y]) //先合计
				num++
			}
		}
	}

	for x := 0; x < 28; x++ {
		for y := 0; y < 28; y++ {
			digits[x][y] = digits[x][y] / float64(num) // 再算平均值
		}
	}

	heatmap := plotter.NewHeatMap(NewHeat(digits), moreland.SmoothBlueRed().Palette(100))

	plt := plot.New()

	plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = 0, 0, 28, 28

	plt.Add(heatmap)
	if err := plt.Save(5*vg.Inch, 5*vg.Inch, "88.png"); err != nil {
		panic(err)
	}
}

// 借鉴代码  https://github.com/clambin/solaredge-monitor/blob/master/plot/grid.go
type Heat struct {
	z   [28][28]float64
	min float64
	max float64
}

func NewHeat(z [28][28]float64) *Heat {
	min := math.Inf(+1)
	max := math.Inf(-1)
	for x := 0; x < 28; x++ {
		for y := 0; y < 28; y++ {
			min = math.Min(min, z[x][y])
			max = math.Max(max, z[x][y])
		}
	}
	return &Heat{
		z:   z,
		min: min,
		max: max,
	}
}

func (g Heat) Dims() (c, r int) {
	return 28, 28
}

func (g Heat) Z(c, r int) float64 {
	return g.z[r][c]
}

func (g Heat) X(c int) float64 {
	return float64(c)
}
func (g Heat) Y(r int) float64 {
	return float64(r)
}

func (g Heat) Min() float64 {
	return g.min
}

func (g Heat) Max() float64 {
	return g.max
}
