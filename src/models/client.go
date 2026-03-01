package models

import (
	"net"
	"sync"

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

	ExpectedKeepalive uint64

	Upstream           *UpstreamServer
	UpstreamConnection net.Conn

	connectionState          int
	connectionStateMutex     sync.RWMutex
	loginPhaseReachedChannel chan bool
}

func (client *DownstreamClient) isAlive() bool {
	return client.connectionState != clientStateKilled
}

func (client *DownstreamClient) IsAlive() bool {
	client.connectionStateMutex.RLock()
	defer client.connectionStateMutex.RUnlock()

	return client.isAlive()
}

func (client *DownstreamClient) IsProxying() bool {
	client.connectionStateMutex.RLock()
	defer client.connectionStateMutex.RUnlock()

	return client.connectionState == clientStateProxying
}

func (client *DownstreamClient) IsInitiating() bool {
	client.connectionStateMutex.RLock()
	defer client.connectionStateMutex.RUnlock()

	return client.connectionState == clientStateInitial
}

func (client *DownstreamClient) Kill() {
	client.connectionStateMutex.Lock()
	defer client.connectionStateMutex.Unlock()

	if client.Upstream != nil {
		client.Upstream.ClientDisconnected(client)
	}

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

func (client *DownstreamClient) RegisterLoginPhaseChannel(loginPhaseReachedChannel chan bool) {
	client.connectionStateMutex.Lock()
	defer client.connectionStateMutex.Unlock()

	if !client.isAlive() {
		return
	}

	client.loginPhaseReachedChannel = loginPhaseReachedChannel
}

func (client *DownstreamClient) LoginFinished() {
	client.connectionStateMutex.RLock()
	defer client.connectionStateMutex.RUnlock()

	if !client.isAlive() {
		return
	}

	if client.loginPhaseReachedChannel != nil {
		go func(finishChan chan bool) {
			finishChan <- true
		}(client.loginPhaseReachedChannel)
	}
}

func (client *DownstreamClient) EnableProxying() {
	client.connectionStateMutex.Lock()
	defer client.connectionStateMutex.Unlock()

	if !client.isAlive() {
		return
	}

	client.connectionState = clientStateProxying
}

func (client *DownstreamClient) StartTransfer() {
	client.connectionStateMutex.Lock()
	defer client.connectionStateMutex.Unlock()

	if !client.isAlive() {
		return
	}

	client.connectionState = clientStateTransferring
}
