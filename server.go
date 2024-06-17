package main

import (
	"net"
	"net/rpc"
	"strconv"
	"sync"
)

type GameState struct {
	Players map[string]*ClientState
	Map     [][]*Elemento
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
		reply.Success = true
	} else {
		reply.Success = false
	}
	return nil
}

func (gs *GameServer) SendCommand(args *CommandArgs, reply *CommandReply) error {
	// TODO
	return nil
}

func (gs *GameServer) GetGameState(args *GameStateArgs, reply *GameStateReply) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	reply.State = []string{gs.state.toString()}
	return nil
}

func start() {
	gameServer := NewGameServer()
	rpc.Register(gameServer)
	listener, err := net.Listen("tcp", ":3696")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}
