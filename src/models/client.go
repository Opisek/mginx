package models

import (
	"net"

	"github.com/google/uuid"
)

const (
	clientStateInitial = iota
	clientStateProxying
	clientStateTransferring
	clientStateKilled
)

type DownstreamClient struct {
	Connection net.Conn

	Version uint64
	Address string
	Port    uint16

	Username string
	Uuid     uuid.UUID

	GamePhase int

	Upstream           *UpstreamServer
	UpstreamConnection net.Conn

	connectionState int
}

func (client *DownstreamClient) IsAlive() bool {
	return client.connectionState != clientStateKilled
}

func (client *DownstreamClient) IsProxying() bool {
	return client.connectionState == clientStateProxying
}

func (client *DownstreamClient) IsInitiating() bool {
	return client.connectionState == clientStateInitial
}

func (client *DownstreamClient) Kill() {
	client.connectionState = clientStateKilled
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

func (client *DownstreamClient) EnableProxying() {
	if !client.IsAlive() {
		return
	}
	client.connectionState = clientStateProxying
}

func (client *DownstreamClient) StartTransfer() {
	if !client.IsAlive() {
		return
	}
	client.connectionState = clientStateTransferring
}
