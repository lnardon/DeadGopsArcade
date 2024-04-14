package main

import (
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
