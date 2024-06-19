package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

func Mover(comando rune) {
	lastMove = comando
	switch comando {
	case 'w':
		playerRef.Move(playerRef.X, playerRef.Y-1, &mapa)
	case 'a':
		playerRef.Move(playerRef.X-1, playerRef.Y, &mapa)
	case 's':
		playerRef.Move(playerRef.X, playerRef.Y+1, &mapa)
	case 'd':
		playerRef.Move(playerRef.X+1, playerRef.Y, &mapa)
	}
}

func adicionaZumbi(x int, y int) {
	zumbi := &Elemento{
		Id:         gerarIdUnico(),
		Tipo:       "zombie",
		Simbolo:    '💀',
		Cor:        termbox.ColorDefault,
		CorFundo:   termbox.ColorDefault,
		Tangivel:   true,
		Interativo: false,
		X:          x,
		Y:          y,
	}
	mapa.AdicionaElemento(zumbi)
	go zumbi.MoverZumbi()
}

func adicionaPlayer(x int, y int) {
	player := &Elemento{
		Id:         gerarIdUnico(),
		Tipo:       "player",
		Simbolo:    '😆',
		Cor:        termbox.ColorDefault,
		CorFundo:   termbox.ColorDefault,
		Tangivel:   true,
		Interativo: false,
		X:          x,
		Y:          y,
	}
	mapa.AdicionaElemento(player)
}

func gerarIdUnico() int {
	mutex.Lock()
	defer mutex.Unlock()

	for {
		id := rand.Int()
		if !idsUsados[id] {
			idsUsados[id] = true
			return id
		}
	}
}

func atirar() {
	tiro := &Elemento{
		Id:         gerarIdUnico(),
		Tipo:       "bullet",
		Simbolo:    '🔹',
		Cor:        termbox.ColorDefault,
		CorFundo:   termbox.ColorDefault,
		Tangivel:   true,
		Interativo: false,
		X:          playerRef.X,
		Y:          playerRef.Y,
	}

	switch lastMove {
	case 'w':
		tiroY = playerRef.Y - 1
		tiroX = playerRef.X
	case 'a':
		tiroY = playerRef.Y
		tiroX = playerRef.X - 1
	case 's':
		tiroY = playerRef.Y + 1
		tiroX = playerRef.X
	case 'd':
		tiroY = playerRef.Y
		tiroX = playerRef.X + 1
	default:
		tiroY = playerRef.Y - 1
		tiroX = playerRef.X
	}

	mapa.AdicionaElemento(tiro)
	tiro.MoveTiro(tiroX, tiroY, &mapa, lastMove)
}

func DropItem(id int) bool {
	var elemento = mapa.GetPositionById(id)

	if true { // 10%
		mapa.RemoveElemento(id)
		item := &Elemento{
			Id:         gerarIdUnico(),
			Tipo:       "item",
			Simbolo:    '💫',
			Cor:        termbox.ColorYellow,
			CorFundo:   termbox.ColorDefault,
			Tangivel:   true,
			Interativo: true,
			X:          elemento.X,
			Y:          elemento.Y,
		}
		mapa.AdicionaElemento(item)
		go func() {
			time.Sleep(2 * time.Second)
			mapa.RemoveElemento(item.Id)
		}()
		return true
	}
	go func() {
		for {
			x := rand.Intn(80)
			y := rand.Intn(30)
			if mapa.GetElemento(x, y).Tipo == "empty" {
				adicionaZumbi(x, y)
				return
			}
		}
	}()
	return false
}

func isWall(x int, y int) bool {
	if mapa.GetElemento(x, y).Tipo == "wall" {
		return true
	}
	return false
}
