package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

type Elemento struct {
	id         int
	simbolo    rune
	cor        termbox.Attribute
	corFundo   termbox.Attribute
	tangivel   bool
	interativo bool
	x          int
	y          int
}

func (element *Elemento) Move(newX int, newY int, mapa *Map) {
	elementoAntigo  := mapa.GetElemento(newX, newY)
    if elementoAntigo.tangivel {
		return
    }
    elementoAntigo.x = element.x
    elementoAntigo.y = element.y
    element.x = newX
	element.y = newY
	mapa.AdicionaElemento(elementoAntigo)
    mapa.RemoveElemento(elementoAntigo.id)
}

func (element *Elemento) MoveTiro(newX int, newY int, mapa *Map, direction rune) {
	if mapa.GetElemento(newX, newY).simbolo == '☠' {
		mapa.RemoveElemento(element.id)
		if dropar(mapa.GetElemento(newX, newY).id) != true {
			mapa.RemoveElemento(mapa.GetElemento(newX, newY).id)
		}
		atualizaMapa()
		return 
	}

	if mapa.GetElemento(newX, newY).tangivel {
		mapa.RemoveElemento(element.id)
		atualizaMapa()
		return
    }
    element.x = newX
    element.y = newY
	atualizaMapa()
	time.Sleep(time.Millisecond * 50)

    switch direction {
    case 'w':
        element.MoveTiro(element.x, element.y-1, mapa, direction)
    case 'a':
        element.MoveTiro(element.x-1, element.y, mapa, direction)
    case 's':
        element.MoveTiro(element.x, element.y+1, mapa, direction)
    case 'd':
        element.MoveTiro(element.x+1, element.y, mapa, direction)
	default:
        element.MoveTiro(element.x, element.y-1, mapa, 'w')
    }
}

func (element *Elemento) SeguePlayer(player *Elemento) {
	for {
		if element.x < player.x {
			element.x++
		} else if element.x > player.x{
			element.x--
		}
		if element.y < player.y {
			element.y++
		} else if element.y > player.y {
			element.y--
		}
		time.Sleep(time.Second)
	}
}

func atualizaMapa() {
	mapa.MontaMapa()
	mapa.DesenhaMapa()
}

func interagir(x int, y int) {
	switch lastMove {
    case 'w':
		if mapa.Mapa[y-1][x].interativo {
			fmt.Println("Interacted and removed")
		}
    case 'a':
		if mapa.Mapa[y][x-1].interativo {
			fmt.Println("Interacted and removed")
		}
    case 's':
		if mapa.Mapa[y+1][x].interativo {
			fmt.Println("Interacted and removed")
		}
    case 'd':
		if mapa.Mapa[y][x+1].interativo {
			fmt.Println("Interacted and removed")
		}
	default:
        return
    }
}