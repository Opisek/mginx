package payloads

import "github.com/google/uuid"

type LoginStart struct {
	Name string
	Uuid uuid.UUID
}

type LoginAcknowledged struct{}

type LoginSuccess struct {
	Name string
	Uuid uuid.UUID
}
