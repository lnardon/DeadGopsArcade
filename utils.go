package main

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

func adicionaZumbi() {
	zumbi := &Elemento{
		id:       gerarIdUnico(),
		simbolo:  '☠',
		cor:      termbox.ColorDefault,
		corFundo: termbox.ColorDefault,
		tangivel: true,
		interativo: false,
		x:        15,
		y:        15,
	}
	mapa.AdicionaElemento(zumbi)
}

func gerarIdUnico() int {
    for {
        id := rand.Int() // gera um número aleatório
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