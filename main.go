package main

import (
	"github.com/nsf/termbox-go"
)

var mapa Map
var playerRef *Elemento
var idsUsados = make(map[int]bool)
var lastMove rune
var tiroX, tiroY int
var maxZombies = 10
var currentZombies = 0

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
			adicionaZumbi()
			currentZombies++
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return // Sair do programa
			}
			if ev.Ch == 'e' {
				interagir()
			} else if ev.Key == termbox.KeySpace {
				go atirar()
			} else {
				Mover(ev.Ch)
			}
			mapa.MontaMapa()
			mapa.DesenhaMapa()	
		}
	}
}