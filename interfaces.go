package main

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
