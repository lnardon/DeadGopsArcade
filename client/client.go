package client

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"

	"github.com/nsf/termbox-go"
)

type GameClient struct {
	server   *rpc.Client
	clientID string
}

func NewGameClient(serverAddress string, clientID string) (*GameClient, error) {
	client, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return &GameClient{
		server:   client,
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

func (gc *GameClient) GetGameState() (GameState, error) {
	args := &GameStateArgs{
		ClientID: gc.clientID,
	}
	reply := &GameStateReply{}
	err := gc.server.Call("GameServer.GetGameState", args, reply)
	if err != nil {
		return GameState{}, err
	}
	return reply.State, nil
}

func (gc *GameClient) GetMap() (Map, error) {
	var mapa Map
	err := gc.server.Call("GameServer.ShowMap", &struct{}{}, &mapa)
	if err != nil {
		fmt.Println("Error to get map:", err)
		return Map{}, err
	}
	return mapa, nil
}

func main() {
	clientID := "exampleClientID"

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What is the server IP + Port to connect?: ")
	serverAddress, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

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

	for {

	}

}
