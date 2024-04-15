package main

import (
	"math/rand"
	"sync"

	"github.com/nsf/termbox-go"
)

var mapa Map
var playerRef *Elemento
var idsUsados = make(map[int]bool)
var lastMove rune
var tiroX, tiroY int
var maxZombies = 10
var currentZombies = 0
var mutex sync.Mutex

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	carregarMapa("map.txt")
	mapa.DesenhaMapa()
	for {
		if currentZombies < maxZombies {
			adicionaZumbi(
				rand.Intn(80),
				rand.Intn(30),
			)
			currentZombies++
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return // exit
			}
			if ev.Ch == 'e' {
				interagir(playerRef.x, playerRef.y)
			}
			
			if ev.Key == termbox.KeySpace {
				go atirar()
			} else {
				Mover(ev.Ch)
			}
			mapa.MontaMapa()
			mapa.DesenhaMapa()	
		}
	}
}