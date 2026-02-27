package parsing

import (
	"errors"
	"fmt"
	util "mginx/protocol/internal"

	"github.com/google/uuid"
)

type LoginStartPayload struct {
	Name string
	Uuid uuid.UUID
}

func ParseLoginStart(buffer []byte) (LoginStartPayload, error) {
	name, buffer, err := util.ParseString(buffer)
	if err != nil {
		return LoginStartPayload{},
			errors.Join(errors.New("could not parse client username"), err)
	}

	uuid, buffer, err := util.ParseUuid(buffer)
	if err != nil {
		return LoginStartPayload{},
			errors.Join(errors.New("could not parse client uuid"), err)
	}

	if len(buffer) != 0 {
		return LoginStartPayload{},
			fmt.Errorf("remaining bytes at the end of payload: %v", len(buffer))
	}

	return LoginStartPayload{
		Name: name,
		Uuid: uuid,
	}, nil
}

type LoginAcknowledgedPayload struct{}

func ParseLoginAcknowledged(buffer []byte) (LoginAcknowledgedPayload, error) {
	if len(buffer) != 0 {
		return LoginAcknowledgedPayload{},
			fmt.Errorf("remaining bytes at the end of payload: %v", len(buffer))
	}

	return LoginAcknowledgedPayload{}, nil
}
