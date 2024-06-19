package main

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"time"

	"github.com/nsf/termbox-go"
)

var numeroComando int = 0
func NewGameClient(serverAddress string, clientID string) (*GameClient, error) {
	client, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return &GameClient{
		server: client,
		clientID: clientID,
	}, nil
}

func (gc *GameClient) Register() (bool, error) {
	args := &RegisterArgs{
		ClientID: gc.clientID,
	}
	fmt.Print(args.ClientID)
	reply := &RegisterReply{}
	err := gc.server.Call("GameServer.RegisterClient", args, reply)
	if err != nil {
		return false, err
	}
	return reply.Success, nil
}

func (gc *GameClient) GetGameState() (GameState, error) {
    args := &GameStateArgs{
        ClientID: gc.clientID,
    }
    reply := &GameStateReply{}
    err := gc.server.Call("GameServer.GetGameState", args, reply)
    if err != nil {
        return GameState{
			Map: nil,
			Players: nil,
		}, err
    }

	//fmt.Println("State:", reply.State.toString())
    return reply.State, nil
}

func (gc *GameClient) GetMap() (*Map, error) {
    args := &ShowMapArgs{}
    reply := &ShowMapReply{}
    err := gc.server.Call("GameServer.ShowMap", args, reply)
    if err != nil {
        fmt.Println("Error to get map:", err)
        return nil, err
    }

    if reply.Map == nil {
        fmt.Println("Received nil map from server")
        return nil, fmt.Errorf("nil map received from server")
    }

	PrintMap(reply.Map)
    return reply.Map, nil
}


func PrintMap(mapa *Map) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, linha := range mapa.Mapa {
		for _, elem := range linha {
			termbox.SetCell(elem.X, elem.Y, elem.Simbolo, elem.Cor, elem.CorFundo)
		}
	}
	termbox.Flush()
}

func main() {
	serverAddress := "localhost:3696"

    randomNumber := rand.Intn(100)
	fmt.Print(string(randomNumber))
	gameClient, err := NewGameClient(serverAddress, string(randomNumber))
    if err != nil {
        fmt.Println("Error to connect in port", err)
        return
    }

	gameClient.Register()

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	tick := time.Tick(250 * time.Millisecond)

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					return 
				}
				if ev.Ch == 'e' {
					 interagir(playerRef.X, playerRef.Y)
				} else if ev.Key == termbox.KeySpace {
					 go atirar()
				} else {
					 Mover(ev.Ch, numeroComando, gameClient)
					 numeroComando++
				}
			}
			_, err := gameClient.GetMap()
			if err != nil {
				fmt.Println("Error to get map:", err)
				return
			}
		case <-tick:
			_, err := gameClient.GetMap()
			if err != nil {
				fmt.Println("Error to get map:", err)
				return
			}

			/*
			state, err := gameClient.GetGameState()
			if err != nil {
				fmt.Println("Error to get game state:", err)
				return
			}
			fmt.Println("State:", state.toString())
			*/
		}
	}
}

