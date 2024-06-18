package main

import (
	"bufio"
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
		mapa.Mapa[elemento.y][elemento.x] = elemento
	}

}

func (mapa *Map) AdicionaIterativa() {
	for _, elemento := range mapa.Elementos {
		if elemento.tipo == "player" || elemento.tipo == "zombie" || elemento.tipo == "bullet" {
			mapa.ThreadsInterativas = append(mapa.ThreadsInterativas, elemento)
		}
	}
}

func (mapa *Map) AdicionaElemento(elemento *Elemento) {
	mutex.Lock()
	defer mutex.Unlock()

	if elemento.tipo == "bullet" && tiroEmExecucao > 5 {
		return
	}

	if elemento.tipo == "bullet" {
		tiroEmExecucao++
	}
	mapa.Elementos = append(mapa.Elementos, elemento)

}

func (mapa *Map) RemoveElemento(id int) {
	for index, elemento := range mapa.Elementos {
		if elemento.id == id {
			vazio := &Elemento{
				id:         gerarIdUnico(),
				tipo:       "empty",
				simbolo:    ' ',
				cor:        termbox.ColorDefault,
				corFundo:   termbox.ColorDefault,
				tangivel:   false,
				interativo: false,
				x:          elemento.x,
				y:          elemento.y,
			}
			mutex.Lock()
			if elemento.tangivel && tiroEmExecucao > 0 {
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
		if elemento.id == id {
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
					id:         gerarIdUnico(),
					tipo:       "zombie",
					simbolo:    '💀',
					cor:        termbox.ColorDefault,
					corFundo:   termbox.ColorDefault,
					tangivel:   true,
					interativo: false,
					x:          x,
					y:          y,
				}
				mapa.AdicionaElemento(zombie)
				break
			case '▤':
				parede := &Elemento{
					id:         gerarIdUnico(),
					tipo:       "wall",
					simbolo:    '▤',
					cor:        termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim,
					corFundo:   termbox.ColorDarkGray,
					tangivel:   true,
					interativo: false,
					x:          x,
					y:          y,
				}
				mapa.AdicionaElemento(parede)
				break
			case '#':
				barreira := &Elemento{
					id:         gerarIdUnico(),
					tipo:       "blockage",
					simbolo:    '🚧',
					cor:        termbox.ColorRed,
					corFundo:   termbox.ColorDefault,
					tangivel:   true,
					interativo: false,
					x:          x,
					y:          y,
				}
				mapa.AdicionaElemento(barreira)
				break
			case '☺':
				personagem := &Elemento{
					id:         gerarIdUnico(),
					tipo:       "player",
					simbolo:    '😆',
					cor:        termbox.ColorBlack,
					corFundo:   termbox.ColorDefault,
					tangivel:   true,
					interativo: false,
					x:          x,
					y:          y,
				}
				playerRef = personagem
				mapa.AdicionaElemento(personagem)
				break
			case ' ':
				vazio := &Elemento{
					id:         gerarIdUnico(),
					tipo:       "empty",
					simbolo:    ' ',
					cor:        termbox.ColorDefault,
					corFundo:   termbox.ColorDefault,
					tangivel:   false,
					interativo: false,
					x:          x,
					y:          y,
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
	aroundElements[0] = mapa.GetElemento(elemento.x, elemento.y+1)
	aroundElements[1] = mapa.GetElemento(elemento.x+1, elemento.y)
	aroundElements[2] = mapa.GetElemento(elemento.x, elemento.y-1)
	aroundElements[3] = mapa.GetElemento(elemento.x-1, elemento.y)
	return aroundElements
}

func SpawnaZumbi() {
	x := rand.Intn(80)
	y := rand.Intn(30)

	if mapa.GetElemento(x, y).tipo == "empty" {
		adicionaZumbi(x, y)
		currentZombies++
		return
	}
	SpawnaZumbi()
}
