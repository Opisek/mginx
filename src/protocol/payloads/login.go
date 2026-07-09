package payloads

import "github.com/google/uuid"

type LoginStart struct {
	Name string
	Uuid uuid.UUID
}

type LoginAcknowledged struct{}

type LoginSuccessOld struct {
	Name string
	Uuid uuid.UUID
}

type LoginSuccess struct {
	Name    string
	Uuid    uuid.UUID
	Session uuid.UUID
}

type LoginDisconnect struct {
	Reason string
}
