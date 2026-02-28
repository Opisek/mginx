package models

import (
	"mginx/config"
	"net"

	"github.com/google/uuid"
)

const (
	StateInitial = iota
	StateProxying
	StateTransferring
	StateKilled
)

type GameClient struct {
	Connection net.Conn

	Version uint64
	Address string
	Port    uint16

	Username string
	Uuid     uuid.UUID

	GamePhase int

	Upstream           *config.ServerConfig
	UpstreamConnection net.Conn

	connectionState int
}

func (client *GameClient) IsAlive() bool {
	return client.connectionState != StateKilled
}

func (client *GameClient) IsProxying() bool {
	return client.connectionState == StateProxying
}

func (client *GameClient) IsInitiating() bool {
	return client.connectionState == StateInitial
}

func (client *GameClient) Kill() {
	client.connectionState = StateKilled
	client.GamePhase = 0xFF

	if client.Connection != nil {
		client.Connection.Close()
		client.Connection = nil
	}

	if client.UpstreamConnection != nil {
		client.UpstreamConnection.Close()
		client.UpstreamConnection = nil
	}
}

func (client *GameClient) EnableProxying() {
	if !client.IsAlive() {
		return
	}
	client.connectionState = StateProxying
}

func (client *GameClient) StartTransfer() {
	if !client.IsAlive() {
		return
	}
	client.connectionState = StateTransferring
}
