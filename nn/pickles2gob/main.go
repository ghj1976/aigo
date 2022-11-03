package main

import (
	"fmt"
	"io"
	"log"
	"os"

	pickle2 "github.com/kisielk/og-rek"
	"github.com/nlpodyssey/gopickle/pickle"
	"github.com/vmihailenco/msgpack"
	"golang.org/x/text/encoding/charmap"
)

func main() {
	Tgopickle()
}

func Tmsgpack() {
	f, err := os.Open("./mnist.pkl")
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	// Latin1是ISO-8859-1的别名，有些环境下写作Latin-1。
	r := charmap.ISO8859_1.NewDecoder().Reader(f)

	buf, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	var out [][][]float32
	err = msgpack.Unmarshal(buf, &out)
	if err != nil {
		panic(err)
	}

	for k, v := range out {
		fmt.Printf("key[%v] value[%v]\n", k, v)
	}
}

func Togrek() {
	f, err := os.Open("./mnist.pkl")
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	// Latin1是ISO-8859-1的别名，有些环境下写作Latin-1。
	r := charmap.ISO8859_1.NewDecoder().Reader(f)

	d := pickle2.NewDecoder(r)
	obj, err := d.Decode()
	if err != nil {
		// panic: Unknown opcode 98 (b) at position 35: 'b'
		log.Panicln(err)
	}

	log.Println(obj)
}

func Tgopickle() {

	// from file
	foo, err := pickle.Load("./mnist.pkl")
	if err != nil {
		log.Panicln(err)
	}
	log.Println(foo)
}
