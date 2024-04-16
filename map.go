package main

import (
	"bufio"
	"os"

	"github.com/nsf/termbox-go"
)

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
		if elemento.simbolo == '☺' || elemento.simbolo == '☠' || elemento.simbolo == '*' {
			mapa.ThreadsInterativas = append(mapa.ThreadsInterativas, elemento)
		}
	}
}

func (mapa *Map) DesenhaMapa() {
	mutex.Lock()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y, linha := range mapa.Mapa {
		for x, elem := range linha {
			termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
		}
	}
	termbox.Flush()
	mutex.Unlock()
}

func (mapa *Map) AdicionaElemento(elemento *Elemento) {
	mutex.Lock()
	defer mutex.Unlock()
	if(elemento.simbolo == '*' && tiroEmExecucao > 5) {
		return
	}

	if(elemento.simbolo == '*') {
		tiroEmExecucao++
	}
	mapa.Elementos = append(mapa.Elementos, elemento)
	
}

func (mapa *Map) RemoveElemento(id int) {
    for index, elemento := range mapa.Elementos {
        if elemento.id == id {
            vazio := &Elemento{
                id:       gerarIdUnico(),
                simbolo:  ' ',
                cor:      termbox.ColorDefault,
                corFundo: termbox.ColorDefault,
                tangivel: false,
				interativo: false,
                x:        elemento.x,
                y:        elemento.y,
            }
			mutex.Lock()
			if(elemento.simbolo == '*' && tiroEmExecucao > 0) {
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
					id: 		gerarIdUnico(),
					simbolo:    '☠',
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
					id: 		gerarIdUnico(),
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
					id: 		gerarIdUnico(),
					simbolo:    '#',
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
					id: 		gerarIdUnico(),
					simbolo:    '☺',
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
					id: 		gerarIdUnico(),
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