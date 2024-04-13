package main

import (
	"time"
	"github.com/nsf/termbox-go"
)

type Map struct {
	Elementos []Elemento
	Mapa [][]Elemento
	ThreadsInterativas []Elemento
}

func (mapa *Map) MontaMapa() {
	if mapa.Mapa == nil {
		mapa.Mapa = make([][]Elemento, 30)
		for i := range mapa.Mapa {
			mapa.Mapa[i] = make([]Elemento, 80)
		}
	}
	
	for _, elemento := range mapa.Elementos {
		mapa.Mapa[elemento.y][elemento.x] = elemento
	}
	
}

func (mapa *Map) AdicionaIterativa() {
	for _, elemento := range mapa.Elementos {
		if(elemento.simbolo == '☺' || elemento.simbolo == '☠' || elemento.simbolo == '*') {
			mapa.ThreadsInterativas = append(mapa.ThreadsInterativas, elemento)
		}
	}
}

func (mapa *Map) DesenhaMapa() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    for y, linha := range mapa.Mapa {
        for x, elem := range linha {
            termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
        }
    }
    termbox.Flush()
}

func (mapa *Map) AdicionaElemento(elemento *Elemento) {
	mapa.Elementos = append(mapa.Elementos, *elemento)
}

func (mapa *Map) RemoveElemento(id int) {
	for index, elemento := range mapa.Elementos {
		if(elemento.id == id) {
			mapa.Elementos = append(mapa.Elementos[:index], mapa.Elementos[index+1:]...)
			vazio := Elemento{
				id: int(time.Nanosecond),
				simbolo: ' ',
				cor: termbox.ColorDefault,
				corFundo: termbox.ColorDefault,
				tangivel: false,
				x: elemento.x,
				y: elemento.y,
			}
			mapa.Elementos[index] = vazio
			return 
		}
	} 
}