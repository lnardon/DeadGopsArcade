package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

var tiroEmExecucao int = 0
var mapa Map = Map{}
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

	for _ = range maxZombies{
		if currentZombies < maxZombies {
			x := rand.Intn(80)
			y := rand.Intn(30)

			if mapa.GetElemento(x, y).tipo == "empty" || mapa.GetElemento(x, y) == nil {
				adicionaZumbi(x, y)
				currentZombies++
			}
		}
	}

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	tick := time.Tick(1000 * time.Millisecond)

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					return // exit
				}
				if ev.Ch == 'e' {
					 interagir(playerRef.x, playerRef.y)
				} else if ev.Key == termbox.KeySpace {
					 go atirar()
				} else {
					 Mover(ev.Ch)
				}
			}
		case <-tick:
			mapa.DesenhaMapa()
		}
	}
}