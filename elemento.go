package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type Elemento struct {
	Id         int
	Simbolo    rune
	Tipo	   string
	Cor        termbox.Attribute
	CorFundo   termbox.Attribute
	Tangivel   bool
	Interativo bool
	X          int
	Y          int
}

func (element *Elemento) Move(newX int, newY int, mapa *Map) {
	elementoAntigo  := mapa.GetElemento(newX, newY)
	if(elementoAntigo.Tipo == "player" && element.Tipo == "zombie" || elementoAntigo.Tipo == "zombie" && element.Tipo == "player") {
		termbox.Close()
		fmt.Println("Game finished! You lose!")
		os.Exit(1)
		return
	}
    if elementoAntigo.Tangivel {
		return
    }
	
    elementoAntigo.X = element.X
    elementoAntigo.Y = element.Y
    element.X = newX
	element.Y = newY
}

func (element *Elemento) MoveTiro(newX int, newY int, mapa *Map, direction rune) {
	if mapa.GetElemento(newX, newY).Tipo == "zombie" {
		
		killedZombies++
		mapa.RemoveElemento(element.Id)
		if DropItem(mapa.GetElemento(newX, newY).Id) != true {
			mapa.RemoveElemento(mapa.GetElemento(newX, newY).Id)
		}
		return 
	}

	if mapa.GetElemento(newX, newY).Tangivel {
		mapa.RemoveElemento(element.Id)
		return
    }
    element.X = newX
    element.Y = newY
	time.Sleep(time.Millisecond * 50)

    switch direction {
    case 'w':
        element.MoveTiro(element.X, element.Y-1, mapa, direction)
    case 'a':
        element.MoveTiro(element.X-1, element.Y, mapa, direction)
    case 's':
        element.MoveTiro(element.X, element.Y+1, mapa, direction)
    case 'd':
        element.MoveTiro(element.X+1, element.Y, mapa, direction)
	default:
        element.MoveTiro(element.X, element.Y-1, mapa, 'w')
    }
}

func interagir(x int, y int) {
	arounds := mapa.GetAround(mapa.GetElemento(x, y));
	for _, elementos :=range arounds {
		if(elementos.Interativo) {
			termbox.Close()
			fmt.Println("Game finished! You won!")
			os.Exit(1)
			return
		}
	}
}

func (el *Elemento) MoverZumbi() {
	for {
		if el.Tipo == "zombie" {
			direction := rune(rand.Intn(4))
			switch direction {
			case 0:
				if !isWall(el.X, el.Y-1){
					el.Move(el.X, el.Y-1, &mapa)
				}
			case 1:
				if !isWall(el.X - 1, el.Y){
					el.Move(el.X-1, el.Y, &mapa)
				}
			
			case 2:
				if !isWall(el.X, el.Y+1){
					el.Move(el.X, el.Y+1, &mapa)
				}
			
			case 3:
				if !isWall(el.X+1, el.Y){
					el.Move(el.X+1, el.Y, &mapa)
				}
			}
		}

		time.Sleep(time.Millisecond * 250)
	}
}