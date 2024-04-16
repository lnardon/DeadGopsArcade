package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type Elemento struct {
	id         int
	simbolo    rune
	tipo	   string
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
	if mapa.GetElemento(newX, newY).tipo == "zombie" {
		mapa.RemoveElemento(element.id)
		if DropItem(mapa.GetElemento(newX, newY).id) != true {
			mapa.RemoveElemento(mapa.GetElemento(newX, newY).id)
		}
		atualizaMapa()
		return 
	}

	if mapa.GetElemento(newX, newY).tangivel {
		mapa.RemoveElemento(element.id)
		mutex.Lock()
		tiroEmExecucao--
		mutex.Unlock()
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

func atualizaMapa() {
	mapa.MontaMapa()
	mapa.DesenhaMapa()
}

func interagir(x int, y int) {
	mutex.Lock()
	defer mutex.Unlock()
	switch lastMove {
    case 'w':
		if mapa.Mapa[y-1][x].interativo {
			termbox.Close()
			fmt.Println("Game finished!")
			os.Exit(1)
		}
    case 'a':
		if mapa.Mapa[y][x-1].interativo {
			termbox.Close()
			fmt.Println("Game finished!")
			os.Exit(1)
		}
    case 's':
		if mapa.Mapa[y+1][x].interativo {
			termbox.Close()
			fmt.Println("Game finished!")
			os.Exit(1)
		}
    case 'd':
		if mapa.Mapa[y][x+1].interativo {
			termbox.Close()
			fmt.Println("Game finished!")
			os.Exit(1)
		}
	default:
        return
    }
}

func (el *Elemento) MoverZumbi() {
	for {
		if el.tipo == "zombie" {
			direction := rune(rand.Intn(4))
			switch direction {
			case 0:
				if !isWall(el.x, el.y-1){
					el.Move(el.x, el.y-1, &mapa)
				}
			case 1:
				if !isWall(el.x - 1, el.y){
					el.Move(el.x-1, el.y, &mapa)
				}
			
			case 2:
				if !isWall(el.x, el.y+1){
					el.Move(el.x, el.y+1, &mapa)
				}
			
			case 3:
				if !isWall(el.x+1, el.y){
					el.Move(el.x+1, el.y, &mapa)
				}
			}
		}

		time.Sleep(time.Millisecond * 250)
		atualizaMapa()
	}
}