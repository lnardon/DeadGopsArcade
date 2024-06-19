package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
)

type GameServer struct {
	clients  map[string]*ClientState
	commands chan CommandArgs
	state    GameState
	mutex    sync.Mutex
}

func NewGameServer() *GameServer {
	return &GameServer{
		clients:  make(map[string]*ClientState),
		commands: make(chan CommandArgs, 100),
		state:    GameState{},
		mutex:    sync.Mutex{},
	}
}

func (gs *GameServer) RegisterClient(args *RegisterArgs, reply *RegisterReply) error {
	fmt.Println("Registering client", args.ClientID)
	position := gs.SpawnClient()
	fmt.Println("Client spawned at", position[0], position[1])

	adicionaPlayer(position[0], position[1], gs.state.Map)

	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if _, exists := gs.clients[args.ClientID]; !exists {
		gs.clients[args.ClientID] = &ClientState{
			ID: 	  args.ClientID,
			PositionX: position[0],
			PositionY: position[1],
			Health:   100,
		}
		reply.Success = true
	} else {
		reply.Success = false
	}

	fmt.Println("Client", gs.clients[args.ClientID].toString(), args.ClientID)
	return nil
}

func (gs *GameServer) SpawnClient() [2]int {
	x := rand.Intn(80)
	y := rand.Intn(30)
	return [2]int{x, y}
}

func (gs *GameServer) SendCommand(args *CommandArgs, reply *CommandReply) error {
	fmt.Println("Received command", args.Command, "from", args.ClientID, args.Command == 'w')
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if client, exists := gs.clients[args.ClientID]; exists {
		switch args.Command {
		case 'w':
			client.PositionY--
			gs.state.Map.Elementos[0].Move(client.PositionX, client.PositionY, gs.state.Map)
		case 'a':
			client.PositionX--
			gs.state.Map.Elementos[0].Move(client.PositionX, client.PositionY, gs.state.Map)
		case 's':
			client.PositionY++
			gs.state.Map.Elementos[0].Move(client.PositionX, client.PositionY, gs.state.Map)
		case 'd':
			client.PositionX++
			gs.state.Map.Elementos[0].Move(client.PositionX, client.PositionY, gs.state.Map)
		}
		reply.Result = "Executed!"
	} else {
		reply.Result = "Error"
	}

	return nil
}

func (gs *GameServer) GetGameState(args *GameStateArgs, reply *GameStateReply) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	reply.State = gs.state
	return nil
}

func (gs *GameServer) ShowMap(args *ShowMapArgs, reply *ShowMapReply) error {
    gs.mutex.Lock()
    defer gs.mutex.Unlock()
    if gs.state.Map == nil {
        fmt.Println("Map is nil")
        return fmt.Errorf("map data is not available")
    }
    reply.Map = gs.state.Map
    return nil
}


func main() {
	carregarMapa("map.txt")
	gameServer := NewGameServer()
	gameServer.state.Map = &mapa
	rpc.Register(gameServer)
	listener, err := net.Listen("tcp", ":3696")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server is waiting to connect on port:", "3696")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		gameServer.RegisterClient(
			&RegisterArgs{
				ClientID: "player",
			},
			&RegisterReply{},
		)
		go rpc.ServeConn(conn)
	}

}
