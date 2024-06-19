package main

import (
	"fmt"
	"net/rpc"
	"strings"

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

func (gc *GameClient) GetMap() (Map, error) {
    args := &ShowMapArgs{}
    reply := &ShowMapReply{}
    err := gc.server.Call("GameServer.ShowMap", args, reply)
    if err != nil {
        fmt.Println("Error to get map:", err)
        return Map{}, err 
    }
    return reply.Map, nil
}

func PrintMap(mapa Map) {
	mutex.Lock()
	defer mutex.Unlock()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y, linha := range mapa.Mapa {
		for x, elem := range linha {
			termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
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

	success, err := gameClient.Register()
	if err != nil {
		fmt.Println("Error to register client", err)
		return
	}
	fmt.Println("Client registered successfully:", success)

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	mapa, err := gameClient.GetMap()
	if err != nil {
		fmt.Println("Error to get map:", err)
		return
	}
	PrintMap(mapa)
}

