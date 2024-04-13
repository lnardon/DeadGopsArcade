package main

import (
    "github.com/nsf/termbox-go"
)

type Elemento struct {
    id int
    simbolo rune
    cor termbox.Attribute
    corFundo termbox.Attribute
    tangivel bool
    interativo bool
    x int
    y int
}


func (element Elemento) Move(newX int, newY int) {
    element.x = newX
	element.y = newY
}