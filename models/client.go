package models

import (
	"net"

	"github.com/google/uuid"
)

type GameClient struct {
	Connection net.Conn
	Version    uint64
	Address    string
	Port       uint16
	Username   string
	Uuid       uuid.UUID
	GamePhase  int
}
