package main

import (
	"net/rpc"
	"strings"
)

type ClientState struct {
	ID 				int
	PositionX       int
	PositionY       int
	Health          int // Maybe?
	LastSequenceNum int
}

type GameClient struct {
	server *rpc.Client
}

func NewGameClient(serverAddress string) (*GameClient, error) {
	client, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return &GameClient{server: client}, nil
}

func (gc *GameClient) Register(clientID string) bool {
	args := &RegisterArgs{ClientID: clientID}
	reply := &RegisterReply{}
	err := gc.server.Call("GameServer.RegisterClient", args, reply)
	return err == nil && reply.Success
}

func (gc *GameClient) SendCommand(command string, sequenceNumber int) string {
	args := &CommandArgs{
		ClientID:       "", // TODO
		Command:        command,
		SequenceNumber: sequenceNumber,
	}
	reply := &CommandReply{}
	gc.server.Call("GameServer.SendCommand", args, reply)
	return reply.Result
}

func (gc *GameClient) GetGameState() string {
	args := &GameStateArgs{ClientID: "exampleClientID"} // TODO
	reply := &GameStateReply{}
	gc.server.Call("GameServer.GetGameState", args, reply)
	return strings.Join(reply.State, ", ")
}

