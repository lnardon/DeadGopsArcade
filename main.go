package main

import (
	"fmt"
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
var maxZombies = 5
var currentZombies = 0
var mutex sync.Mutex
var killedZombiesMsg string
var killedZombies int

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	carregarMapa("map.txt")

	for _ = range maxZombies{
		SpawnaZumbi()
	}

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	tick := time.Tick(100* time.Millisecond)

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
			AtualizaMapa()
		}
	}
}

func SpawnaZumbi() {
	x := rand.Intn(80)
	y := rand.Intn(30)

	if mapa.GetElemento(x, y).tipo == "empty"  {
		adicionaZumbi(x, y)
		currentZombies++
		return
	}
	SpawnaZumbi()
}

func desenhaBarraDeStatus() {
	killedZombiesMsg = fmt.Sprintf("Você matou %d zombies", killedZombies)
    for i, c := range killedZombiesMsg {
        termbox.SetCell(i, len(mapa.Mapa)+1, c, termbox.ColorBlack, termbox.ColorDefault)
    }
    msg := "Use WASD para mover e E para interagir. ESC para sair."
    for i, c := range msg {
        termbox.SetCell(i, len(mapa.Mapa)+5, c, termbox.ColorBlack, termbox.ColorDefault)
    }

	msg1 := "Ao matar o zumbi você terá 1 segundo para interagir com o item e ganhar o jogo."
    for i, c := range msg1 {
        termbox.SetCell(i, len(mapa.Mapa)+2, c, termbox.ColorBlack, termbox.ColorDefault)
    }
	msg2 := "Se o zumbi chegar em você, você perde."
    for i, c := range msg2 {
        termbox.SetCell(i, len(mapa.Mapa)+3, c, termbox.ColorBlack, termbox.ColorDefault)
    }
}

