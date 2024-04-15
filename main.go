package main

import (
	"bufio"
	"os"
	"math/rand"

	"github.com/nsf/termbox-go"
)

var mapa Map

var efeitoNeblina = false
var revelado [][]bool
var raioVisao int = 3
var playerRef *Elemento
var idsUsados = make(map[int]bool)
var lastMove rune
var tiroX, tiroY int

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	carregarMapa("map.txt")
	mapa.DesenhaMapa()
	adicionaZumbi()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return // Sair do programa
			}
			if ev.Ch == 'e' {
				interagir()
			} else if ev.Key == termbox.KeySpace {
				go atirar()
			} else {
				mover(ev.Ch)
			}
			mapa.MontaMapa()
			mapa.DesenhaMapa()	
		}
	}
}

func interagir() {
	//
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

func mover(comando rune) {
	lastMove = comando
	switch comando {
	case 'w':
		playerRef.Move(playerRef.x, playerRef.y-1, &mapa)
	case 'a':
		playerRef.Move(playerRef.x-1, playerRef.y, &mapa)
	case 's':
		playerRef.Move(playerRef.x, playerRef.y+1, &mapa)
	case 'd':
		playerRef.Move(playerRef.x+1, playerRef.y, &mapa)
	}
}

// adiciona zumbi para teste, pois se ele já estiver no mapa não da para matar
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
