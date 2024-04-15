package main

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

func Mover(comando rune) {
	lastMove = comando
	switch comando {
	case 'w':
		playerRef.Move(playerRef.x, playerRef.y-1, &mapa)
	case 'a':
		playerRef.Move(playerRef.x-1, playerRef.y, &mapa)
	case 's':
		playerRef.Move(playerRef.x, playerRef.y+1, &mapa)
	case 'd':
		playerRef.Move(playerRef.x+1, playerRef.y, &mapa)
	}
}

func adicionaZumbi(x int, y int) {
	zumbi := &Elemento{
		id:       gerarIdUnico(),
		simbolo:  '☠',
		cor:      termbox.ColorDefault,
		corFundo: termbox.ColorDefault,
		tangivel: true,
		interativo: false,
		x:        x,
		y:        y,
	}
	mapa.AdicionaElemento(zumbi)
}

func gerarIdUnico() int {
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
        id:       gerarIdUnico(),
        simbolo:  '*',
        cor:      termbox.ColorDefault,
        corFundo: termbox.ColorDefault,
        tangivel: false,
        interativo: false,
        x:        playerRef.x,
        y:        playerRef.y,
    }

	switch lastMove {
	case 'w':
		tiroY = playerRef.y - 1
		tiroX = playerRef.x
	case 'a':
		tiroY = playerRef.y
		tiroX = playerRef.x - 1
	case 's':
		tiroY = playerRef.y + 1
		tiroX = playerRef.x
	case 'd':
		tiroY = playerRef.y
		tiroX = playerRef.x + 1
	default:
		tiroY = playerRef.y - 1
		tiroX = playerRef.x
	}

    mapa.AdicionaElemento(tiro)
    tiro.MoveTiro(tiroX, tiroY, &mapa, lastMove)
}