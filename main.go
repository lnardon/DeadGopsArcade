package main

import (
	"bufio"
	"os"

	"github.com/nsf/termbox-go"
)

var mapa Map

var efeitoNeblina = false
var revelado [][]bool
var raioVisao int = 3
var playerRef *Elemento

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	carregarMapa("map.txt")
	mapa.DesenhaMapa()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return // Sair do programa
			}
			if ev.Ch == 'e' {
				interagir()
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
			case '♣':
				vegetacao := &Elemento{
					simbolo:    '♣',
					cor:        termbox.ColorGreen,
					corFundo:   termbox.ColorDefault,
					tangivel:   false,
					interativo: false,
					x:          x,
					y:          y,
				}
				mapa.AdicionaElemento(vegetacao)
				break
			case '☺':
				personagem := &Elemento{
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
