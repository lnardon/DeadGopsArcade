package server

import (
	"sync"

	"github.com/nsf/termbox-go"
)

type RegisterArgs struct {
	ClientID string
}

type RegisterReply struct {
	Success bool
}

type CommandArgs struct {
	ClientID       string
	Command        string
	SequenceNumber int
}

type CommandReply struct {
	Result string
}

type GameStateArgs struct {
	ClientID string
}

type GameStateReply struct {
	State GameState
}

type GameServerInterface interface {
	RegisterClient(args *RegisterArgs, reply *RegisterReply) error
	SendCommand(args *CommandArgs, reply *CommandReply) error
	GetGameState(args *GameStateArgs, reply *GameStateReply) error
}

type GameState struct {
	Players map[string]*ClientState
	Map     Map
}

type GameServer struct {
	clients  map[string]*ClientState
	commands chan CommandArgs
	state    GameState
	mutex    sync.Mutex
}

type Map struct {
	Elementos          []*Elemento
	Mapa               [][]*Elemento
	ThreadsInterativas []*Elemento
}

type Elemento struct {
	id         int
	simbolo    rune
	tipo       string
	cor        termbox.Attribute
	corFundo   termbox.Attribute
	tangivel   bool
	interativo bool
	x          int
	y          int
}

type ClientState struct {
	ID              string
	PositionX       int
	PositionY       int
	Health          int
	LastSequenceNum int
}
