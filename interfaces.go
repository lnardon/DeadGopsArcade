package main

import (
	"fmt"
	"net/rpc"
	"strings"
)

type RegisterArgs struct {
	ClientID int
}

type RegisterReply struct {
	Success bool
}

type CommandArgs struct {
	ClientID       int
	Command        rune
	SequenceNumber int
}

type CommandReply struct {
	Result string
}

type GameClient struct {
	server *rpc.Client
	clientID int
}

type GameStateArgs struct {
	ClientID int
}

type GameState struct {
	Players map[string]*ClientState
	Map     *Map
}

type GameStateReply struct {
	State GameState
}

type ShowMapArgs struct{}

type ShowMapReply struct {
    Map *Map
}

type GameServerInterface interface {
	RegisterClient(args *RegisterArgs, reply *RegisterReply) error
	SendCommand(args *CommandArgs, reply *CommandReply) error
	GetGameState(args *GameStateArgs, reply *GameStateReply) error
}

func (gc *GameClient) SendCommand(command rune, sequenceNumber int) (string, error) {
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

func (cs *ClientState) toString() string {
	return fmt.Sprintf("ID: %s, Position: (%d, %d), Health: %d, LastSeqNum: %d",
		cs.ID, cs.PositionX, cs.PositionY, cs.Health, cs.LastSequenceNum)
}

func (e *Elemento) toString() string {
	return fmt.Sprintf("Elemento ID: %d, Type: %s, Symbol: %c, Position: (%d, %d), Tangible: %t, Interactive: %t",
		e.Id, e.Tipo, e.Simbolo, e.X, e.Y, e.Tangivel, e.Interativo)
}

func (gs *GameState) toString() string {
	var sb strings.Builder
	sb.WriteString("Game State:\n")
	for id, player := range gs.Players {
		sb.WriteString(fmt.Sprintf("Player %s: %s\n", id, player.toString()))
	}
	sb.WriteString(gs.Map.toString())
	return sb.String()
}