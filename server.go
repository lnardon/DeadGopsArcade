package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"strconv"
	"sync"
)

type GameState struct {
	Players map[string]*ClientState
	Map     *Map
}

func (gs *GameState) toString() string {
	var state string
	for id, player := range gs.Players {
		state += id + " - Pos: (" + strconv.Itoa(player.PositionX) + ", " + strconv.Itoa(player.PositionY) + "), Health: " + strconv.Itoa(player.Health) + "\n"
	}
	return state
}

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
		gs.clients[args.ClientID] = &ClientState{}
		gs.clients[args.ClientID].PositionX = position[0]
		gs.clients[args.ClientID].PositionY = position[1]
		reply.Success = true
	} else {
		reply.Success = false
	}
	return nil
}

func (gs *GameServer) SpawnClient() [2]int {
	x := rand.Intn(80)
	y := rand.Intn(30)
	return [2]int{x, y}
}

func (gs *GameServer) SendCommand(args *CommandArgs, reply *CommandReply) error {
	// Command processing logic goes here
	return nil
}

func (gs *GameServer) GetGameState(args *GameStateArgs, reply *GameStateReply) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	reply.State = []string{gs.state.toString()}
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
	// for _ = range maxZombies {
	// 	SpawnaZumbi()
	// }

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
