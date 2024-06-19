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
	Map     Map
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
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if _, exists := gs.clients[args.ClientID]; !exists {
		gs.clients[args.ClientID] = &ClientState{}
		position := gs.SpawnClient()
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

	if gs.state.Map.GetElemento(x, y).tipo == "empty" {
		adicionaPlayer(x, y)
		return [2]int{x, y}
	}
	gs.SpawnClient()
	return [2]int{0, 0}
}

func (gs *GameServer) SendCommand(args *CommandArgs, reply *CommandReply) error {
	if(args.SequenceNumber == gs.clients[args.ClientID].LastSequenceNum) {
		reply.Result = "command already processed"
		return nil
	}
	return nil
}

func (gs *GameServer) GetGameState(args *GameStateArgs, reply *GameStateReply) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	reply.State = gs.state
	return nil
}

// Nao sei se ta certo esse metodo
func (gs *GameServer) ShowMap() Map {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.state.Map
}

func main() {
	carregarMapa("map.txt")
	//for _ = range maxZombies {
	//	SpawnaZumbi()
	//}

	gameServer := NewGameServer()
	gameServer.state.Map = mapa
	rpc.Register(gameServer)
	listener, err := net.Listen("tcp", ":3696")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		fmt.Println("Server is waiting to connect in port:", "3696")
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}