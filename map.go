package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"

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

type Map struct {
	Elementos          []*Elemento
	Mapa               [][]*Elemento
	ThreadsInterativas []*Elemento
}

func (mapa *Map) MontaMapa() {
	if mapa.Mapa == nil {
		mapa.Mapa = make([][]*Elemento, 30)
		for i := range mapa.Mapa {
			mapa.Mapa[i] = make([]*Elemento, 80)
		}
	}

	for _, elemento := range mapa.Elementos {
		mapa.Mapa[elemento.Y][elemento.X] = elemento
	}

}

func (mapa *Map) AdicionaIterativa() {
	for _, elemento := range mapa.Elementos {
		if elemento.Tipo == "player" || elemento.Tipo == "zombie" || elemento.Tipo == "bullet" {
			mapa.ThreadsInterativas = append(mapa.ThreadsInterativas, elemento)
		}
	}
}

func (mapa *Map) AdicionaElemento(elemento *Elemento) {
    mutex.Lock()
    defer mutex.Unlock()

    if elemento.Tipo == "bullet" && tiroEmExecucao > 5 {
        return
    }

    if elemento.Tipo == "bullet" {
        tiroEmExecucao++
    }

    mapa.Elementos = append(mapa.Elementos, elemento)

	if mapa.Mapa != nil {
		mapa.MontaMapa()
	}
    fmt.Println(mapa.toString())
}


func (mapa *Map) RemoveElemento(id int) {
	for index, elemento := range mapa.Elementos {
		if elemento.Id == id {
			vazio := &Elemento{
				Id:         gerarIdUnico(),
				Tipo:       "empty",
				Simbolo:    ' ',
				Cor:        termbox.ColorDefault,
				CorFundo:   termbox.ColorDefault,
				Tangivel:   false,
				Interativo: false,
				X:          elemento.X,
				Y:          elemento.Y,
			}
			mutex.Lock()
			if elemento.Tangivel && tiroEmExecucao > 0 {
				tiroEmExecucao--
			}
			mapa.Elementos[index] = vazio
			mapa.Elementos = append(mapa.Elementos[:index], mapa.Elementos[index+1:]...)
			mutex.Unlock()
			return
		}
	}
}

func (mapa *Map) GetElemento(x int, y int) *Elemento {
	mutex.Lock()
	defer mutex.Unlock()
	return mapa.Mapa[y][x]
}

func (mapa *Map) GetPositionById(id int) *Elemento {
	for _, elemento := range mapa.Elementos {
		if elemento.Id == id {
			return elemento
		}
	}
	return nil
}

func carregarMapa(nomeArquivo string) {
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		panic(err)
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	y := 0
	x := 0
	for scanner.Scan() {
		linhaTexto := scanner.Text()
		for _, char := range linhaTexto {
			switch char {
			case '☠':
				zombie := &Elemento{
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
				mapa.AdicionaElemento(zombie)
				break
			case '▤':
				parede := &Elemento{
					Id:         gerarIdUnico(),
					Tipo:       "wall",
					Simbolo:    '▤',
					Cor:        termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim,
					CorFundo:   termbox.ColorDarkGray,
					Tangivel:   true,
					Interativo: false,
					X:          x,
					Y:          y,
				}
				mapa.AdicionaElemento(parede)
				break
			case '#':
				barreira := &Elemento{
					Id:         gerarIdUnico(),
					Tipo:       "blockage",
					Simbolo:    '🚧',
					Cor:        termbox.ColorRed,
					CorFundo:   termbox.ColorDefault,
					Tangivel:   true,
					Interativo: false,
					X:          x,
					Y:          y,
				}
				mapa.AdicionaElemento(barreira)
				break
			case '☺':
				personagem := &Elemento{
					Id:         gerarIdUnico(),
					Tipo:       "player",
					Simbolo:    '😆',
					Cor:        termbox.ColorBlack,
					CorFundo:   termbox.ColorDefault,
					Tangivel:   true,
					Interativo: false,
					X:          x,
					Y:          y,
				}
				playerRef = personagem
				mapa.AdicionaElemento(personagem)
				break
			case ' ':
				vazio := &Elemento{
					Id:         gerarIdUnico(),
					Tipo:       "empty",
					Simbolo:    ' ',
					Cor:        termbox.ColorDefault,
					CorFundo:   termbox.ColorDefault,
					Tangivel:   false,
					Interativo: false,
					X:          x,
					Y:          y,
				}
				mapa.AdicionaElemento(vazio)
				break

			}
			x++
		}
		x = 0
		y++
	}
	mapa.AdicionaIterativa()
	mapa.MontaMapa()
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (mapa *Map) GetAround(elemento *Elemento) []*Elemento {
	aroundElements := make([]*Elemento, 4)
	aroundElements[0] = mapa.GetElemento(elemento.X, elemento.Y+1)
	aroundElements[1] = mapa.GetElemento(elemento.X+1, elemento.Y)
	aroundElements[2] = mapa.GetElemento(elemento.X, elemento.Y-1)
	aroundElements[3] = mapa.GetElemento(elemento.X-1, elemento.Y)
	return aroundElements
}

func SpawnaZumbi() {
	x := rand.Intn(80)
	y := rand.Intn(30)

	if mapa.GetElemento(x, y).Tipo == "empty" {
		adicionaZumbi(x, y)
		currentZombies++
		return
	}
	SpawnaZumbi()
}

func(mapa *Map)toString() string {
	str := ""
	for _, linha := range mapa.Mapa {
		for _, elem := range linha {
			str += string(elem.Simbolo)
		}
		str += "\n"
	}
	return str
}