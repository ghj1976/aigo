package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// 给棋盘的不同动作生成zobrist哈希
func main() {
	rand.Seed(time.Now().Unix())

	buf := strings.Builder{}
	buf.WriteString("package aigo\n")
	buf.WriteString("\n")

	buf.WriteString("var (\n")
	buf.WriteString(fmt.Sprintf("	EmptyBoardHashCode = int64(%d)\n", rand.Int63n(math.MaxInt64)))
	buf.WriteString("	BoardPointHashCode = MapInit()\n")
	buf.WriteString(")\n")
	buf.WriteString("\n")

	buf.WriteString("func MapInit() map[BoardPoint]int64 {\n")
	buf.WriteString("	mappp := make(map[BoardPoint]int64)\n")
	for i := 1; i <= 19; i++ {
		for j := 1; j <= 19; j++ {
			buf.WriteString(fmt.Sprintf("	mappp[NewBoardPoint(%d, %d, None)] = int64(%d)\n", i, j, rand.Int63n(math.MaxInt64)))
			buf.WriteString(fmt.Sprintf("	mappp[NewBoardPoint(%d, %d, Black)] = int64(%d)\n", i, j, rand.Int63n(math.MaxInt64)))
			buf.WriteString(fmt.Sprintf("	mappp[NewBoardPoint(%d, %d, White)] = int64(%d)\n", i, j, rand.Int63n(math.MaxInt64)))
		}
	}

	buf.WriteString("	return mappp\n")
	buf.WriteString("}\n")
	buf.WriteString("\n")

	WriteString2File("../../zobrist.go", buf.String())
}

// WriteString2File 写文件
func WriteString2File(shortFileName, txt string) {
	filename := path.Join(".", shortFileName)
	log.Println(filename)

	fo, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer fo.Close()
	fo.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM ， 避免文件乱码
	fo.WriteString(txt)

	if err != nil {
		log.Println(err)
		return
	}
}
