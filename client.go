package main

import (
	"fmt"
	"net/rpc"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

type GameClient struct {
	server *rpc.Client
	clientID string
}

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
	reply := &RegisterReply{}
	err := gc.server.Call("GameServer.RegisterClient", args, reply)
	if err != nil {
		return false, err
	}
	return reply.Success, nil
}

func (gc *GameClient) SendCommand(command string, sequenceNumber int) (string, error) {
    args := &CommandArgs{
        ClientID:       gc.clientID,
        Command:        command,
        SequenceNumber: sequenceNumber,
    }
    reply := &CommandReply{}
    err := gc.server.Call("GameServer.SendCommand", args, reply)
    if err != nil {
        return "", err
    }
    return reply.Result, nil
}

func (gc *GameClient) GetGameState() (string, error) {
    args := &GameStateArgs{
        ClientID: gc.clientID,
    }
    reply := &GameStateReply{}
    err := gc.server.Call("GameServer.GetGameState", args, reply)
    if err != nil {
        return "", err
    }
    return strings.Join(reply.State, ", "), nil
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

    fmt.Println(reply.Map.toString(), "mapa")
    return reply.Map, nil
}


func PrintMap(mapa Map) {
	mutex.Lock()
	defer mutex.Unlock()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y, linha := range mapa.Mapa {
		for x, elem := range linha {

			termbox.SetCell(x, y, elem.Simbolo, elem.Cor, elem.CorFundo)
		}
	}
	termbox.Flush()
}


func main() {
	serverAddress := "localhost:3696"
	clientID := "exampleClientID"

	gameClient, err := NewGameClient(serverAddress, clientID)
    if err != nil {
        fmt.Println("Error to connect in port", err)
        return
    }

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

	tick := time.Tick(250* time.Millisecond)

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
					 Mover(ev.Ch)
				}
			}
		case <-tick:
			_, err := gameClient.GetMap()
			if err != nil {
				fmt.Println("Error to get map:", err)
				return
			}
			PrintMap(mapa)
			
		}
	}
}

